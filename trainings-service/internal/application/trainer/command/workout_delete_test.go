package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/mocks"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldDeleteWorkoutGroupOwnedByTrainerWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().DeleteWorkoutGroup(ctx, workout.UUID()).Return(nil)

	SUT := command.NewWorkoutDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, workout.UUID(), workout.TrainerUUID())

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := command.NewWorkoutDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Equal(command.ErrWorkoutGroupNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().DeleteWorkoutGroup(ctx, workout.UUID()).Return(expectedErr)

	SUT := command.NewWorkoutDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, workout.UUID(), workout.TrainerUUID())

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}
