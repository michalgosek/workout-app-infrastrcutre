package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldScheduleWorkoutToSpecifiedWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup

	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	workoutGroupWithCustomer.AssignCustomer(customerDetails)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, workoutGroupWithCustomer).Return(nil)
	customerRepository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).
		Run(func(ctx context.Context, workout customer.WorkoutDay) {
			assertions.Equal(workoutGroup.UUID(), workout.GroupUUID())
			assertions.Equal(workoutGroup.Date(), workout.Date())
			assertions.Equal(customerUUID, workout.CustomerUUID())
		}).Return(nil)

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerDetails.UUID(),
		CustomerName: customerDetails.Name(),
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Nil(err)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenQueryWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const groupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerName: "John Doe",
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const customerName = "John Doe"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup
	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	workoutGroupWithCustomer.AssignCustomer(customerDetails)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, workoutGroupWithCustomer).Return(errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:

	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerName: customerDetails.Name(),
		CustomerUUID: customerDetails.UUID(),
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenUpsertCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup

	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	workoutGroupWithCustomer.AssignCustomer(customerDetails)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, workoutGroupWithCustomer).Return(nil)
	customerRepository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout customer.WorkoutDay) {
		assertions.Equal(workoutGroup.UUID(), workout.GroupUUID())
		assertions.Equal(workoutGroup.Date(), workout.Date())
		assertions.Equal(customerUUID, workout.CustomerUUID())
	}).Return(errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:

	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerName: customerDetails.Name(),
		CustomerUUID: customerDetails.UUID(),
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenWorkoutGroupDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const workoutGroupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutGroupUUID).Return(trainer.WorkoutGroup{}, nil)
	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: "John Doe",
		GroupUUID:    workoutGroupUUID,
	})

	// then:
	assertions.Equal(err, command.ErrResourceNotFound)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}
