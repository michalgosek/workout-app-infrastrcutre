package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnTrainerGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("68580811-9ebd-452e-8f4d-948e1de59646", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("8f35342d-bc59-4ac2-815a-c9c2477cd71f", trainer, date)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &training)

	SUT := query.NewTrainerGroupHandler(cli, query.Config{
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

	expectedGroup := createExpectedTrainerGroup(cli, training.UUID())

	// when:
	trainingGroup, err := SUT.Do(ctx, training.UUID(), trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, trainingGroup)
}
