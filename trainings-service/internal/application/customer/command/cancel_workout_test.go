package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldCancelWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	repository := new(mocks.CancelWorkoutHandlerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(repository)

	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	trainerWorkout.AssignCustomer(customerUUID)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(nil)

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenQueryWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	repository := new(mocks.CancelWorkoutHandlerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(repository)

	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const groupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	repository.EXPECT().QueryWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	repository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenDeleteCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	repository := new(mocks.CancelWorkoutHandlerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(repository)

	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	trainerWorkout.AssignCustomer(customerUUID)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	repository.AssertExpectations(t)
}

func TestCancelWorkoutHandlerShouldReturnErrorWhenUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	repository := new(mocks.CancelWorkoutHandlerRepository)
	ctx := context.Background()
	SUT := command.NewCancelWorkoutHandler(repository)

	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	trainerWorkout.AssignCustomer(customerUUID)

	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(errors.New("err"))

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
	})

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	repository.AssertExpectations(t)
}
