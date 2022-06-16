package command_test

import (
	"context"
	"errors"
	"testing"

	command2 "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/mocks"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldNotAssignCustomerToWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := command2.NewAssignCustomerHandler(repository)

	// when:
	err := SUT.Do(ctx, command2.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workoutUUID,
	})

	// then:
	assertions.Equal(application.ErrScheduleNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToSpecifiedWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutdWithCustomer := workout
	workoutdWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutdWithCustomer).Return(expectedErr)

	SUT := command2.NewAssignCustomerHandler(repository)

	// when:
	err := SUT.Do(ctx, command2.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workout.UUID(),
	})

	// then:
	assertions.ErrorIs(err, command2.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldAssignCustomerToSpecifiedWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutWithCustomer := workout
	_ = workoutWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutWithCustomer).Return(nil)

	SUT := command2.NewAssignCustomerHandler(repository)

	// when:
	err := SUT.Do(ctx, command2.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workout.UUID(),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}
