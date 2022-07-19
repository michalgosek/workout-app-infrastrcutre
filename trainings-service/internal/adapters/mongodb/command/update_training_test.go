package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldUpdateTrainingWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("860da5fd-3346-4b71-b579-b5daf8897f05", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	}
	insertTrainingHandler := command.NewInsertTrainingHandler(cli, commandCfg)
	SUT := command.NewUpdateTrainingHandler(cli, commandCfg)

	defer func() {
		db := cli.Database("insert_training_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	_ = insertTrainingHandler.Do(ctx, &trainingGroup)
	_ = trainingGroup.AssignParticipant(newTestParticipant("c6975a21-a098-4c94-a7de-de01a731b57a"))

	// when:
	err := SUT.Do(ctx, &trainingGroup)

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, trainingGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup := documents.UnmarshalToTrainingGroup(writeModel)
	assertions.Equal(trainingGroup, actualWorkoutDomainGroup)
}
