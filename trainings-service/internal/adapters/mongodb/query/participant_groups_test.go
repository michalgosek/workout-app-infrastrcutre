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
	firstTrainer := testutil.NewTestTrainer("7b4ebc81-7b88-42d1-825e-7af6c8251e55", "John Doe1")
	secondTrainer := testutil.NewTestTrainer("5d639d22-b1d8-4c2e-9140-ec43f6fe4395", "Jerry Smith2")
	participant := testutil.NewTestParticipant("1d48af67-38f0-4832-83f5-a090368534e9", "Chris Barker")

	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("3d599126-f202-48d5-a860-7c58f9832e3f", firstTrainer, date)
	_ = firstTraining.AssignParticipant(participant)

	secondTraining := testutil.NewTestTrainingGroup("44eefe17-aac2-43a3-9fc8-5e17032867e8", secondTrainer, date)
	_ = secondTraining.AssignParticipant(participant)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 10 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &firstTraining)
	_ = handler.InsertTrainingGroup(ctx, &secondTraining)

	SUT := query.NewParticipantGroupsHandler(cli, query.Config{
		Database:     testutil.DatabaseName,
		Collection:   testutil.CollectionName,
		QueryTimeout: 10 * time.Second,
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
