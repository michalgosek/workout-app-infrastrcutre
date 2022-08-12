package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestShouldDeleteTrainingGroupsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("77de3c0e-ec98-45bd-99c3-c6ef0524fb06", "John Doe")
	date := newTestStaticTime()
	firstTrainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)
	secondTrainingGroup := newTestTrainingGroup("d65ede26-89e7-436e-b9d8-0142844a3905", trainer, date)

	cli := newTestMongoClient()
	commandCfg := command.Config{
		Database:       DatabaseName,
		Collection:     CollectionName,
		CommandTimeout: 5 * time.Second,
	}
	handler := command.NewInsertTrainingGroupHandler(cli, commandCfg)
	_ = handler.InsertTrainingGroup(ctx, &firstTrainingGroup)
	_ = handler.InsertTrainingGroup(ctx, &secondTrainingGroup)

	SUT := command.NewDeleteTrainingGroupsHandler(cli, commandCfg)
	defer func() {
		db := cli.Database(DatabaseName)
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.DeleteTrainingGroups(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := findTrainingGroup(cli, firstTrainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)

	writeModel, err = findTrainingGroup(cli, secondTrainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}
