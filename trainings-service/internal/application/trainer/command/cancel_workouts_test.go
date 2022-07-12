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

func TestShouldCancelTrainerWorkoutsGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	group := testutil.NewTrainerWorkoutGroup(trainerUUID)
	service := mocks.NewTrainingsService(t)

	service.EXPECT().CancelTrainerWorkoutGroups(ctx, trainerUUID).Return(nil)

	SUT, _ := command.NewCancelWorkoutsHandler(service)

	// when:
	err := SUT.Do(ctx, group.TrainerUUID())

	// then:
	assertions.Nil(err)
	service.AssertExpectations(t)
}

func TestShouldNotCancelTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	ctx := context.Background()
	service := mocks.NewTrainingsService(t)
	repositoryFailureErr := errors.New("repository failure")
	service.EXPECT().CancelTrainerWorkoutGroups(ctx, trainerUUID).Return(repositoryFailureErr)

	SUT, _ := command.NewCancelWorkoutsHandler(service)

	// when:
	err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	service.AssertExpectations(t)
}
