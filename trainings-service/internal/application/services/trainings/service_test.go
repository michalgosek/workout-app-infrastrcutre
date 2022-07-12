package trainings_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestService_ShouldCancelCustomerWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
		groupUUID    = "438a1155-0145-4087-9ff8-16ebae0877d3"
		trainerUUID  = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	)
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	ctx := context.Background()
	SUT, _ := trainings.NewService(customerService, trainerService)

	customerService.EXPECT().CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	}).Return(nil)

	trainerService.EXPECT().CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(nil)

	// when:
	err := SUT.CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, trainerService, customerService)
}

func TestService_ShouldNotCancelCustomerWorkoutDayWhenCustomerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
		groupUUID    = "438a1155-0145-4087-9ff8-16ebae0877d3"
		trainerUUID  = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	)
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	ctx := context.Background()
	SUT, _ := trainings.NewService(customerService, trainerService)

	repositoryErr := errors.New("repository failure")
	customerService.EXPECT().CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	}).Return(repositoryErr)

	// when:
	err := SUT.CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryErr)
	mock.AssertExpectationsForObjects(t, trainerService, customerService)
}

func TestService_ShouldNotCancelCustomerWorkoutDayWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
		groupUUID    = "438a1155-0145-4087-9ff8-16ebae0877d3"
		trainerUUID  = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	)

	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	ctx := context.Background()
	SUT, _ := trainings.NewService(customerService, trainerService)

	customerService.EXPECT().CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	}).Return(nil)

	repositoryErr := errors.New("repository failure")
	trainerService.EXPECT().CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(repositoryErr)

	// when:
	err := SUT.CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryErr)
	mock.AssertExpectationsForObjects(t, trainerService, customerService)
}

func TestService_ShouldAssignCustomerToWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	repositoryFailureErr := errors.New("repository failure")
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(trainer.AssignedCustomerWorkoutGroupDetails{}, repositoryFailureErr)

	// when:
	err := SUT.AssignCustomerToWorkoutGroup(ctx, trainings.AssignCustomerToWorkoutArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestService_ShouldNotAssignCustomerToWorkoutGroupWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	repositoryFailureErr := errors.New("repository failure")
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(trainer.AssignedCustomerWorkoutGroupDetails{}, repositoryFailureErr)

	// when:
	err := SUT.AssignCustomerToWorkoutGroup(ctx, trainings.AssignCustomerToWorkoutArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestService_ShouldNotAssignCustomerToWorkoutGroupWhenCustomerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		customerName = "John Doe"
		groupUUID    = "e233ef39-37df-492b-8736-a106b6f14363"
		trainerUUID  = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	assignedWorkoutDetails := trainer.AssignedCustomerWorkoutGroupDetails{
		UUID:        groupUUID,
		TrainerUUID: trainerUUID,
		Name:        "dummy",
		Date:        time.Now().AddDate(0, 0, 1),
	}
	trainerService.EXPECT().AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(assignedWorkoutDetails, nil)

	repositoryFailureErr := errors.New("repository failure")
	customerService.EXPECT().ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
		Date:         assignedWorkoutDetails.Date,
	}).Return(repositoryFailureErr)

	// when:
	err := SUT.AssignCustomerToWorkoutGroup(ctx, trainings.AssignCustomerToWorkoutArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestService_ShouldCancelTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		groupUUID   = "e233ef39-37df-492b-8736-a106b6f14363"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	trainerService.EXPECT().CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	}).Return(nil)

	customerService.EXPECT().CancelWorkoutDaysWithGroup(ctx, groupUUID).Return(nil)

	// when:
	err := SUT.CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestService_ShouldNotCancelTrainerWorkoutGroupWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		groupUUID   = "e233ef39-37df-492b-8736-a106b6f14363"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	repositoryFailure := errors.New("repository failure")
	trainerService.EXPECT().CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	}).Return(repositoryFailure)

	// when:
	err := SUT.CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailure)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}

func TestService_ShouldNotCancelTrainerWorkoutGroupWhenCustomerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "504868b4-89b8-48bc-9da3-213d90f0c91e"
		groupUUID   = "e233ef39-37df-492b-8736-a106b6f14363"
	)

	ctx := context.Background()
	customerService := new(mocks.CustomerService)
	trainerService := new(mocks.TrainerService)
	SUT, _ := trainings.NewService(customerService, trainerService)

	trainerService.EXPECT().CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	}).Return(nil)
	repositoryFailure := errors.New("repository failure")
	customerService.EXPECT().CancelWorkoutDaysWithGroup(ctx, groupUUID).Return(repositoryFailure)

	// when:
	err := SUT.CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailure)
	mock.AssertExpectationsForObjects(t, customerService, trainerService)
}
