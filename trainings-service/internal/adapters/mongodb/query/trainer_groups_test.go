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
	trainer := testutil.NewTestTrainer("b367f443-2591-491a-9b5d-46304654df26", "Jane Doe")
	date := testutil.NewTestStaticTime()
	firstTraining := testutil.NewTestTrainingGroup("a505c5e2-d886-4d62-92df-d37fa6d8d0ea", trainer, date)
	secondTraining := testutil.NewTestTrainingGroup("6d5d8c41-53af-4bc6-8b1b-76b2fa48dc69", trainer, date)

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
