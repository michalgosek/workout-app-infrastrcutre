package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldUpdateTrainingGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("29de02f6-6398-4a5a-b450-4600fff3e271", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("c5d78199-79af-4f95-957b-f6806a3edb2b", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	}
	handler := command.NewInsertTrainingGroupHandler(cli, commandCfg)
	SUT := command.NewUpdateTrainingGroupHandler(cli, commandCfg)

	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	_ = handler.InsertTrainingGroup(ctx, &training)
	_ = training.AssignParticipant(newTestParticipant("4af091ed-fd0f-4fa9-bcc9-e57989e6a458"))

	// when:
	err := SUT.UpdateTrainingGroup(ctx, &training)

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, training.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	expectedGroup := documents.UnmarshalToTrainingGroup(writeModel)
	assertions.Equal(training, expectedGroup)
}
