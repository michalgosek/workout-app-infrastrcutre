package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestShouldDeleteTrainerGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	firstTrainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)
	secondTrainingGroup := newTestTrainingGroup("d65ede26-89e7-436e-b9d8-0142844a3905", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	}
	insertTrainingHandler := command.NewInsertTrainerGroupHandler(cli, commandCfg)
	_ = insertTrainingHandler.Do(ctx, &firstTrainingGroup)
	_ = insertTrainingHandler.Do(ctx, &secondTrainingGroup)

	SUT := command.NewDeleteTrainerGroupsHandler(cli, commandCfg)
	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.Do(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, firstTrainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)

	writeModel, err = findTrainingGroup(cli, secondTrainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}
