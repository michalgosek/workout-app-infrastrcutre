package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnAllTrainingGroupsReadModelsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	firstTraining := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)
	secondTraining := newTestTrainingGroup("619a0b4d-5509-40bf-b7ff-704f15adc406", trainer, date)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainingHandler(cli, command.Config{
		Database:       "insert_training_db",
		Collection:     "trainings",
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &firstTraining)
	_ = insertTrainingHandler.Do(ctx, &secondTraining)

	SUT := query.NewAllTrainingsHandler(cli, query.Config{
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

	expectedReadModels, err := createAllTrainingGroupReadModels(cli)
	assertions.Nil(err)

	// when:
	trainings, err := SUT.Do(ctx)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedReadModels, trainings)
}
