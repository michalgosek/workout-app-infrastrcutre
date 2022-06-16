package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/mocks"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldDeleteWorkoutsGroupOwnedByTrainerWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().DeleteWorkoutGroups(ctx, trainerUUID).Return(nil)

	SUT := command.NewWorkoutsDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, workout.TrainerUUID())

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteWorkoutsGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	var workouts []trainer.WorkoutGroup
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(workouts, nil)
	SUT := command.NewWorkoutsDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Equal(application.ErrScheduleNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteTrainerWorkoutsGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().DeleteWorkoutGroups(ctx, trainerUUID).Return(expectedErr)
	SUT := command.NewWorkoutsDeleteHandler(repository)

	// when:
	err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Contains(err.Error(), expectedErr.Error())
	repository.AssertExpectations(t)
}
