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

func TestShouldReturnTrainingGroupWithSpecifiedUUIDSuccessfully_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("978b4c99-caa8-4909-8cee-5828b17b7a9e", "John Doe")
	date := testutil.NewTestStaticTime()
	expected := testutil.NewTestTrainingGroup("e2e17d29-c6e9-4a2c-8d44-ff131f63a614", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &expected)

	SUT := query.NewTrainingGroupHandler(cli, query.Config{
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

	// when:
	actual, err := SUT.TrainingGroup(ctx, expected.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expected, actual)
}
