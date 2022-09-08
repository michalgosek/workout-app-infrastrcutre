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
	trainer := newTestTrainer("6fd96ff0-6ad7-4e8e-a1a9-dd270cc8278e", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("ba86a74e-bce8-4ac1-a430-a81d70862827", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		CommandTimeout: 10 * time.Second,
		Collection:     CollectionName,
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
	_ = training.AssignParticipant(newTestParticipant("7137c664-2767-4591-88b7-ef2ca4c2354a"))

	// when:
	err := SUT.UpdateTrainingGroup(ctx, &training)

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, training.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	expectedGroup := documents.ConvertToTrainingGroup(writeModel)
	assertions.Equal(training, expectedGroup)
}
