package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldInsertTrainerGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	cli := newTestMongoClient()
	SUT := command.NewInsertTrainerGroupHandler(cli, command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database("insert_training_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

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
