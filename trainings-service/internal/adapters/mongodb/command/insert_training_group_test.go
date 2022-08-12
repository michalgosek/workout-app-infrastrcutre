package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldInsertTrainingGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("f9e836f3-d2e3-4d37-9a7e-d7b6d1856e7a", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("40dbf5d6-e223-4f05-a8a8-531e3c4cb634", trainer, date)

	cli := newTestMongoClient()

	SUT := command.NewInsertTrainingGroupHandler(cli, command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.InsertTrainingGroup(ctx, &training)

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, training.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	expectedGroup := documents.ConvertToTrainingGroup(writeModel)
	assertions.Equal(training, expectedGroup)
}
