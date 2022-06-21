package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldCancelWorkoutGroupOwnedByTrainerWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	group := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, group.UUID()).Return(group, nil)
	repository.EXPECT().DeleteTrainerWorkoutGroup(ctx, group.UUID()).Return(nil)
	SUT := command.NewCancelWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		GroupUUID:   group.UUID(),
		TrainerUUID: group.TrainerUUID(),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotCancelWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	group := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutUUID).Return(group, nil)
	SUT := command.NewCancelWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		GroupUUID:   workoutUUID,
		TrainerUUID: trainerUUID,
	})

	// then:
	assertions.Equal(command.ErrWorkoutGroupNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotCancelTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	group := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)

	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, group.UUID()).Return(group, nil)
	repository.EXPECT().DeleteTrainerWorkoutGroup(ctx, group.UUID()).Return(expectedErr)
	SUT := command.NewCancelWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.CancelWorkout{
		GroupUUID:   group.UUID(),
		TrainerUUID: group.TrainerUUID(),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}
