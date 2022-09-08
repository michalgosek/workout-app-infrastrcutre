package notifications_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/notifications"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/notifications/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/authorization"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

func TestServiceShouldCreateNotificationWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	httpCli := mocks.NewHTTPClient(t)
	SUT := notifications.NewService(httpCli)
	ctx := context.WithValue(context.Background(), authorization.Token, "69b9d020-57e7-4fa4-9b02-517fd85d803a")
	cmd := command.Notification{
		UserUUID:     "eae2c856-7559-4675-8413-8ef6d6d27099",
		TrainingUUID: "ef54063a-2274-4c54-be76-226f6d331755",
		Title:        "training name",
		Content:      "training desc",
		Trainer:      "John Doe",
		Date:         time.Time{},
	}

	httpResponse := &http.Response{StatusCode: http.StatusCreated}
	httpCli.EXPECT().Do(mock.Anything).Return(httpResponse, nil)

	// when:
	err := SUT.CreateNotification(ctx, cmd)

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, httpCli)
}
