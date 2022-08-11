package command_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"notification-service/internal/adapters/mongodb/client"
	"notification-service/internal/adapters/mongodb/command"
	"notification-service/internal/adapters/mongodb/testutil"
	"notification-service/internal/domain"
	"testing"
	"time"
)

func TestShouldInsertUserNotificationWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const userUUID = "d0995436-edc0-48cd-b8a8-0c48ff56afed"
	const trainingUUID = "5a530775-b7f6-4019-8398-8a014d57ca5c"

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
	SUT := command.NewInsertNotificationHandler(cli, command.InsertNotificationHandlerConfig{
		Database:       testCliCfg.Database,
		CommandTimeout: testCliCfg.CommandTimeout,
	})
	expected := testutil.NewTestUserNotification(userUUID, trainingUUID)

	// when:
	err := SUT.InsertNotification(ctx, expected)

	// then:
	assertions.Nil(err)

	doc, err := testCli.FindInsertedNotification(userUUID)
	assertions.Nil(err)
	assertions.NotEmpty(doc)

	actual := domain.ConvertToDomainNotification(doc)
	assertions.Nil(err)
	assertions.Equal(expected, actual)
}
