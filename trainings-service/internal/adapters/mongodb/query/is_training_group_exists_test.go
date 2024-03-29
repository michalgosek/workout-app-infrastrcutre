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

func TestShouldReturnTrueWhenFoundTrainingGroup_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("b17d1854-ed44-48cd-b69b-ca7586dede0b", "John Doe")
	date := testutil.NewTestStaticTime()
	training := testutil.NewTestTrainingGroup("966e3b51-f9ed-4637-ab03-ca56d82d5a5e", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &training)

	SUT := query.NewDuplicateTrainingGroupHandler(cli, query.Config{
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
	exists, err := SUT.IsTrainingGroupExists(ctx, &training)

	// then:
	assertions.Nil(err)
	assertions.True(exists)
}
