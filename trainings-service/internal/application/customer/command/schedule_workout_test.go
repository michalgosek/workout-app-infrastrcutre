package command_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerDetails.UUID(), workoutGroup.UUID()).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutGroup.TrainerUUID(), workoutGroup.UUID()).Return(workoutGroup, nil)
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
		TrainerUUID:  workoutGroup.TrainerUUID(),
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, trainerRepository, customerRepository)
}

func TestShouldReturnErrorWhenQueryWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const groupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerName: "John Doe",
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	mock.AssertExpectationsForObjects(t, trainerRepository, customerRepository)
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

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, workoutGroup.UUID()).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, workoutGroup.UUID()).Return(workoutGroup, nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, workoutGroupWithCustomer).Return(errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerName: customerDetails.Name(),
		CustomerUUID: customerDetails.UUID(),
		GroupUUID:    workoutGroup.UUID(),
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	mock.AssertExpectationsForObjects(t, trainerRepository, customerRepository)
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

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, workoutGroup.UUID()).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, workoutGroup.UUID()).Return(workoutGroup, nil)
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
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	mock.AssertExpectationsForObjects(t, customerRepository, trainerRepository)
}

func TestShouldReturnErrorWhenWorkoutGroupDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const workoutGroupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, workoutGroupUUID).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, workoutGroupUUID).Return(trainer.WorkoutGroup{}, nil)
	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: "John Doe",
		GroupUUID:    workoutGroupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Equal(err, command.ErrResourceNotFound)
	mock.AssertExpectationsForObjects(t, customerRepository, trainerRepository)
}

func TestShouldReturnErrorWhenAttemptsToScheduleDuplicateWorkoutGroup_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	workoutDay := testutil.NewWorkoutDay(customerUUID)
	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, workoutDay.GroupUUID()).Return(workoutDay, nil)

	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    workoutDay.GroupUUID(),
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(command.ErrWorkoutGroupDuplicated, err)
	mock.AssertExpectationsForObjects(t, trainerRepository, customerRepository)
}

func TestShouldReturnErrorWhenAttemptsToScheduleWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const groupUUID = "2e286992-114f-4cdd-a9e7-d4ee8ef37ea9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	SUT := command.NewScheduleWorkoutHandler(customerRepository, trainerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, command.ErrResourceNotFound)
	mock.AssertExpectationsForObjects(t, customerRepository, trainerRepository)
}
