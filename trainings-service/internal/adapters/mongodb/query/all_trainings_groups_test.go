package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnAllTrainingGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("4de6329a-85e9-4b1e-8b39-d9e1c64f6d7c", "John Doe")
	date := newTestStaticTime()
	firstTraining := newTestTrainingGroup("84a40429-c8da-4eac-a693-67bf79456a2a", trainer, date)
	secondTraining := newTestTrainingGroup("297cad1d-9b5b-41eb-b505-a30992d9746f", trainer, date)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &firstTraining)
	_ = insertTrainingHandler.Do(ctx, &secondTraining)

	SUT := query.NewAllTrainingsHandler(cli, query.Config{
		Database:     DatabaseName,
		Collection:   CollectionName,
		QueryTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	expectedGroups, err := createExpectedAllTrainingGroups(cli)
	assertions.Nil(err)

	// when:
	trainingGroups, err := SUT.Do(ctx)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups, trainingGroups)
}
