package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	rm "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldReturnTrainerGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("ad35908e-771f-4f89-93d6-4edb749e4c5d", "John Doe")
	date := newTestStaticTime()
	firstTraining := newTestTrainingGroup("dd1fb0a1-4fd5-41c6-b79d-2a34f23b8f2c", trainer, date)
	secondTraining := newTestTrainingGroup("b05dba7d-ffa2-401e-86b3-b40babc0ab21", trainer, date)

	cli := newTestMongoClient()
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &firstTraining)
	_ = insertTrainingHandler.Do(ctx, &secondTraining)

	SUT := query.NewTrainerGroupsHandler(cli, query.Config{
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

	expectedGroups := []rm.TrainerGroup{
		createExpectedTrainerGroup(cli, firstTraining.UUID()),
		createExpectedTrainerGroup(cli, secondTraining.UUID()),
	}

	// when:
	trainings, err := SUT.Do(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups, trainings)
}
