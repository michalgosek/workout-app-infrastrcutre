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

func TestShouldReturnAllTrainingGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("eb4f831a-3093-424e-a3be-b522369ab8ed", "John Doe")
	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("966ec5be-e82a-41f2-a7fa-0d88f702c6f5", trainer, date)
	secondTraining := testutil.NewTestTrainingGroup("bfe123db-42f7-4d97-8f38-0bea73891e09", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &firstTraining)
	_ = handler.InsertTrainingGroup(ctx, &secondTraining)

	SUT := query.NewAllTrainingsHandler(cli, query.Config{
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

	expectedGroups, err := testutil.CreateAllTrainingGroups(cli)
	assertions.Nil(err)

	// when:
	trainingGroups, err := SUT.AllTrainingGroups(ctx)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups, trainingGroups)
}
