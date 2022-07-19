package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestShouldDeleteTrainingWithSuccess(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	}
	insertTrainingHandler := command.NewInsertTrainingHandler(cli, commandCfg)
	_ = insertTrainingHandler.Do(ctx, &trainingGroup)

	SUT := command.NewDeleteTrainingsHandler(cli, commandCfg)
	defer func() {
		db := cli.Database("insert_training_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.Do(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, trainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}
