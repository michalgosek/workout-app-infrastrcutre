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
	trainer := testutil.NewTestTrainer("eb86bba3-2e42-494d-ba7e-bc638a04a526", "Mark Doe")
	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("02d66236-a647-4088-9c88-a3e2be3cb979", trainer, date)
	secondTraining := testutil.NewTestTrainingGroup("65556967-5f61-44fc-812c-3a42fb848abb", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 10 * time.Second,
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
