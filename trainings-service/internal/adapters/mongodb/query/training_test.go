package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnTrainingReadModelWithSuccess(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainingHandler(cli, command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &trainingGroup)

	SUT := query.NewTrainingHandler(cli, query.Config{
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

	expectedReadModel := createTrainingGroupReadModel(cli, trainingGroup.UUID())

	// when:
	training, err := SUT.Do(ctx, trainingGroup.UUID(), trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedReadModel, training)
}
