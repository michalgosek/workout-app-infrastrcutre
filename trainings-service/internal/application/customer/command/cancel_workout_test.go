package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCancelWorkoutHandler_ShouldCancelCustomerWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
		groupUUID    = "438a1155-0145-4087-9ff8-16ebae0877d3"
		trainerUUID  = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	)
	trainingsService := new(mocks.TrainingsService)

	ctx := context.Background()
	SUT, _ := command.NewCancelWorkoutHandler(trainingsService)

	trainingsService.EXPECT().CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, trainingsService)
}
