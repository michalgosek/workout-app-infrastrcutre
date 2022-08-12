package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"

	application "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnAllTrainerParticipantsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("eb4f831a-3093-424e-a3be-b522369ab8ed", "John Doe")

	firstParticipant := testutil.NewTestParticipant("8663d6d3-a373-49dc-9a01-8c4a37e69445", "Jerry Kayne")
	secondParticipant := testutil.NewTestParticipant("d32fe5c0-d06d-445f-b822-ab1dd73ee2ad", "Jane Williams")

	date := testutil.NewTestStaticTime()

	firstTraining := testutil.NewTestTrainingGroup("3a1b0c0e-76dd-4349-b9cb-d30c0bad4be9", trainer, date)
	_ = firstTraining.AssignParticipant(firstParticipant)

	secondTraining := testutil.NewTestTrainingGroup("0022b37f-15cb-4076-8b2c-4c0d61e8c44a", trainer, date)
	_ = secondTraining.AssignParticipants(secondParticipant)

	expectedParticipants := []application.TrainerParticipant{
		{
			TrainingUUID: firstTraining.UUID(),
			TrainingName: firstTraining.Name(),
			Trainer:      trainer.Name(),
			UserUUID:     firstParticipant.UUID(),
			Date:         firstTraining.Date(),
		},
		{
			TrainingUUID: secondTraining.UUID(),
			TrainingName: secondTraining.Name(),
			Trainer:      trainer.Name(),
			UserUUID:     secondParticipant.UUID(),
			Date:         secondTraining.Date(),
		},
	}

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &firstTraining)
	_ = handler.InsertTrainingGroup(ctx, &secondTraining)

	SUT := query.NewTrainerParticipantsHandler(cli, query.Config{
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

	// when:
	participants, err := SUT.TrainerParticipants(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedParticipants, participants)
}
