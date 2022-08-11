package query_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"notification-service/internal/adapters/mongodb/client"
	"notification-service/internal/adapters/mongodb/command"
	"notification-service/internal/adapters/mongodb/query"
	"notification-service/internal/adapters/mongodb/testutil"
	application "notification-service/internal/application/query"
	"testing"
	"time"
)

func TestShouldGetUserNotificationsWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const userUUID = "119507d9-e9f7-4166-b421-ed2304ff3cb0"
	const firstTrainingUUID = "71063e6a-dc84-4ccd-b7f9-d80d052a4e15"
	const secondTrainingUUID = "0ee780db-9054-40d0-8aea-2e904d2800e0"

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

	expected := []application.Notification{
		{
			UUID:         first.UUID(),
			UserUUID:     first.UserUUID(),
			TrainingUUID: first.TrainingUUID(),
			Title:        first.Title(),
			Trainer:      first.Trainer(),
			Content:      first.Content(),
			Date:         first.Date().Format(application.UIFormat),
		},
		{
			UUID:         second.UUID(),
			UserUUID:     second.UserUUID(),
			TrainingUUID: second.TrainingUUID(),
			Title:        second.Title(),
			Trainer:      second.Trainer(),
			Content:      second.Content(),
			Date:         second.Date().Format(application.UIFormat),
		},
	}

	_ = inserter.InsertNotification(ctx, first)
	_ = inserter.InsertNotification(ctx, second)

	SUT := query.NewAllNotificationHandler(cli, query.AllNotificationHandlerConfig{
		Database:       testCliCfg.Database,
		CommandTimeout: testCliCfg.CommandTimeout,
	})

	// when:
	actual, err := SUT.AllNotifications(ctx, userUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expected, actual)
}
