package command_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"notification-service/internal/adapters/mongodb/client"
	"notification-service/internal/adapters/mongodb/command"
	"notification-service/internal/adapters/mongodb/testutil"
	"testing"
	"time"
)

func TestShouldClearUserNotificationsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const userUUID = "d4a237f3-2c79-4284-a836-a30586959782"
	const firstTrainingUUID = "f5e8fffc-615f-41e3-b3cb-744e9d7f733e"
	const secondTrainingUUID = "a02698e4-1b93-45cf-a973-b79d96f4654e"

	ctx := context.Background()
	testCliCfg := testutil.TestClientConfig{
		Addr:           "mongodb://localhost:27017",
		CommandTimeout: 10 * time.Second,
		QueryTimeout:   10 * time.Second,
		Database:       "notification_service_test_db",
		Collection:     command.NewUserCollectionName(userUUID),
	}
	testCli := testutil.NewTestMongoClient(testCliCfg)

	defer func() {
		err := testCli.Drop(ctx)
		if err != nil {
			panic(err)
		}
		err = testCli.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	cli, _ := client.New(testCliCfg.Addr, 5*time.Second)
	inserter := command.NewInsertNotificationHandler(cli, command.InsertNotificationHandlerConfig{
		Database:       testCliCfg.Database,
		CommandTimeout: testCliCfg.CommandTimeout,
	})
	first := testutil.NewTestUserNotification(userUUID, firstTrainingUUID)
	second := testutil.NewTestUserNotification(userUUID, secondTrainingUUID)

	_ = inserter.InsertNotification(ctx, first)
	_ = inserter.InsertNotification(ctx, second)

	SUT := command.NewClearNotificationsHandler(cli, command.ClearNotificationsHandlerConfig{
		Database:       testCliCfg.Database,
		CommandTimeout: testCliCfg.CommandTimeout,
	})

	// when:
	err := SUT.ClearNotifications(ctx, userUUID)

	// then:
	assertions.Nil(err)

	docs, err := testCli.FindInsertedNotifications(userUUID)
	assertions.Nil(err)
	assertions.Empty(docs)
}
