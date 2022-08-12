package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnTrainerGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("68580811-9ebd-452e-8f4d-948e1de59646", "John Doe")
	date := testutil.NewTestStaticTime()
	training := testutil.NewTestTrainingGroup("8f35342d-bc59-4ac2-815a-c9c2477cd71f", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &training)

	SUT := query.NewTrainerGroupHandler(cli, query.Config{
		Database:     testutil.DatabaseName,
		Collection:   testutil.CollectionName,
		QueryTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database(testutil.DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	expectedGroup := testutil.CreateTrainerGroup(cli, training.UUID())

	// when:
	trainingGroup, err := SUT.TrainerGroup(ctx, training.UUID(), trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, trainingGroup)
}
