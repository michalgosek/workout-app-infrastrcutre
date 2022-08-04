package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnParticipantGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	firstTrainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	secondTrainer := newTestTrainer("f6d61286-e4fd-4922-9f53-1c96164c4920", "Jerry Smith")
	participant := newTestParticipant("f5e13337-0842-4af0-917b-fe2a14c4325b", "Chris Barker")

	date := newTestStaticTime()
	firstTraining := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", firstTrainer, date)
	_ = firstTraining.AssignParticipant(participant)

	secondTraining := newTestTrainingGroup("619a0b4d-5509-40bf-b7ff-704f15adc406", secondTrainer, date)
	_ = secondTraining.AssignParticipant(participant)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &firstTraining)
	_ = insertTrainingHandler.Do(ctx, &secondTraining)

	SUT := query.NewParticipantGroupsHandler(cli, query.Config{
		Database:     "insert_training_db",
		Collection:   "trainings",
		QueryTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database("insert_training_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	expectedParticipantGroups, err := createExpectedParticipantTrainingGroups(cli, participant.UUID())
	assertions.Nil(err)

	// when:
	participantGroups, err := SUT.Do(ctx, participant.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedParticipantGroups, participantGroups)
}
