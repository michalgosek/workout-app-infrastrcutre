package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnParticipantGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	firstTrainer := testutil.NewTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	secondTrainer := testutil.NewTestTrainer("f6d61286-e4fd-4922-9f53-1c96164c4920", "Jerry Smith")
	participant := testutil.NewTestParticipant("f5e13337-0842-4af0-917b-fe2a14c4325b", "Chris Barker")

	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", firstTrainer, date)
	_ = firstTraining.AssignParticipant(participant)

	secondTraining := testutil.NewTestTrainingGroup("619a0b4d-5509-40bf-b7ff-704f15adc406", secondTrainer, date)
	_ = secondTraining.AssignParticipant(participant)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &firstTraining)
	_ = handler.InsertTrainingGroup(ctx, &secondTraining)

	SUT := query.NewParticipantGroupsHandler(cli, query.Config{
		Database:     testutil.DatabaseName,
		Collection:   testutil.CollectionName,
		QueryTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database(testutil.DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	expectedParticipantGroups, err := testutil.CreateParticipantTrainingGroups(cli, participant.UUID())
	assertions.Nil(err)

	// when:
	participantGroups, err := SUT.ParticipantGroups(ctx, participant.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedParticipantGroups, participantGroups)
}
