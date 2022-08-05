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
	trainer := newTestTrainer("eb4f831a-3093-424e-a3be-b522369ab8ed", "John Doe")
	date := newTestStaticTime()
	firstTraining := newTestTrainingGroup("966ec5be-e82a-41f2-a7fa-0d88f702c6f5", trainer, date)
	secondTraining := newTestTrainingGroup("bfe123db-42f7-4d97-8f38-0bea73891e09", trainer, date)

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
