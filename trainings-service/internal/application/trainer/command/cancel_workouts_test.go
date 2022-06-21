package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldCancelWorkoutsGroupsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	group := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().DeleteTrainerWorkoutGroups(ctx, trainerUUID).Return(nil)

	SUT := command.NewCancelWorkoutsHandler(repository)

	// when:
	err := SUT.Do(ctx, group.TrainerUUID())

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotCancelWorkoutsGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().DeleteTrainerWorkoutGroups(ctx, trainerUUID).Return(expectedErr)
	SUT := command.NewCancelWorkoutsHandler(repository)

	// when:
	err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Equal(command.ErrRepositoryFailure, err)
	repository.AssertExpectations(t)
}
