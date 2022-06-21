package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldCancelWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const customerName = "John Doe"

	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(customerRepository, trainerRepository)

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(customerDetails)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(nil)

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Nil(err)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenQueryWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const groupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	SUT := command.NewCancelWorkoutHandler(customerRepository, trainerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenDeleteCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	const customerName = "Jonn Doe"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	SUT := command.NewCancelWorkoutHandler(customerRepository, trainerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(customerDetails)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(customerRepository, trainerRepository)

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(customerDetails)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}
