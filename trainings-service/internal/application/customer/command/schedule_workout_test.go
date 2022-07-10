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

func TestScheduleWorkoutHandler_ShouldScheduleCustomerWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	ctx := context.Background()
	trainingsService := new(mocks.TrainingsService)
	SUT, _ := command.NewScheduleWorkoutHandler(trainingsService)

	trainingsService.EXPECT().AssignCustomerToWorkoutGroup(ctx, trainings.AssignCustomerToWorkoutArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkoutArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, trainingsService)
}
