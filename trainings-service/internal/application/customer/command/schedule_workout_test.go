package command_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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

	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)

	ctx := context.Background()
	SUT, _ := command.NewScheduleWorkoutHandler(trainerService, customerService)

	assignedWorkoutDetails := services.AssignedCustomerWorkoutGroupDetails{
		UUID:        groupUUID,
		TrainerUUID: trainerUUID,
		Name:        "dummy",
		Date:        time.Now().AddDate(0, 0, 1),
	}
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(assignedWorkoutDetails, nil)

	customerService.EXPECT().ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         assignedWorkoutDetails.Date,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestScheduleWorkoutHandler_ShouldNotScheduleCustomerWorkoutDayWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)

	ctx := context.Background()
	SUT, _ := command.NewScheduleWorkoutHandler(trainerService, customerService)

	repositoryFailureErr := errors.New("repository failure")
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(services.AssignedCustomerWorkoutGroupDetails{}, repositoryFailureErr)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestScheduleWorkoutHandler_ShouldNotScheduleCustomerWorkoutDayWhenCustomerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)

	ctx := context.Background()
	SUT, _ := command.NewScheduleWorkoutHandler(trainerService, customerService)

	assignedWorkoutDetails := services.AssignedCustomerWorkoutGroupDetails{
		UUID:        groupUUID,
		TrainerUUID: trainerUUID,
		Name:        "dummy",
		Date:        time.Now().AddDate(0, 0, 1),
	}
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(assignedWorkoutDetails, nil)

	repositoryFailureErr := errors.New("repository failure")
	customerService.EXPECT().ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         assignedWorkoutDetails.Date,
	}).Return(repositoryFailureErr)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}
