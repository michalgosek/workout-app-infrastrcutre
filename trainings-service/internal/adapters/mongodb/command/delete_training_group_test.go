package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestShouldDeleteTrainingGroupWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	training := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	}
	handler := command.NewInsertTrainingGroupHandler(cli, commandCfg)
	_ = handler.InsertTrainingGroup(ctx, &training)

	SUT := command.NewDeleteTrainingGroupHandler(cli, commandCfg)
	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
		err = cli.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.DeleteTrainingGroup(ctx, training.UUID(), trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, training.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}
