package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query/testutil"
	rm "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnTrainerGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := testutil.NewTestTrainer("2fb39fe8-4fec-4312-afe6-113de4f3360f", "John Doe")
	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("86059f1e-95c8-4666-8440-fbd9572c147c", trainer, date)
	secondTraining := testutil.NewTestTrainingGroup("293ea857-1fac-4074-ad95-83f78c2ce112", trainer, date)

	cli := testutil.NewTestMongoClient()
	handler := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       testutil.DatabaseName,
		Collection:     testutil.CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = handler.InsertTrainingGroup(ctx, &firstTraining)
	_ = handler.InsertTrainingGroup(ctx, &secondTraining)

	SUT := query.NewTrainerGroupsHandler(cli, query.Config{
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

	expectedGroups := []rm.TrainerGroup{
		testutil.CreateTrainerGroup(cli, firstTraining.UUID()),
		testutil.CreateTrainerGroup(cli, secondTraining.UUID()),
	}

	// when:
	trainings, err := SUT.TrainerGroups(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups, trainings)
}
