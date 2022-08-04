package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldUpdateTrainerGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("efa9510c-f3ef-4942-9f6f-c71ecccd4790", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("cd91677d-5aaa-404c-ab5c-3218f92fa580", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	}
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, commandCfg)
	SUT := command.NewUpdateTrainerGroupHandler(cli, commandCfg)

	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	_ = insertTrainingHandler.Do(ctx, &training)
	_ = training.AssignParticipant(newTestParticipant("4af091ed-fd0f-4fa9-bcc9-e57989e6a458"))

	// when:
	err := SUT.Do(ctx, &training)

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, training.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	expectedGroup := documents.UnmarshalToTrainingGroup(writeModel)
	assertions.Equal(training, expectedGroup)
}
