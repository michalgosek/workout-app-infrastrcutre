package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldGetRequestedTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	ctx := context.Background()
	group := testutil.NewTrainerWorkoutGroup(trainerUUID)
	expectedGroups := query.ConvertToWorkoutGroupsDetails(group)

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, group.UUID()).Return(group, nil)
	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	actualGroup, err := SUT.Do(ctx, trainerUUID, group.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups[0], actualGroup)
	repository.AssertExpectations(t)
}

func TestShouldGetEmptyTrainerWorkoutGroupWhenRequestedGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const groupUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	group := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(group, nil)

	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	actualGroup, err := SUT.Do(ctx, trainerUUID, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
	repository.AssertExpectations(t)
}

func TestShouldNotGetTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const groupUUID = "ef547a4e-f0ef-4282-a308-e985cbac2a01"
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, expectedError)
	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	workoutGroupDetails, err := SUT.Do(ctx, trainerUUID, groupUUID)

	// then:
	assertions.ErrorContains(err, err.Error())
	assertions.Empty(workoutGroupDetails)
	repository.AssertExpectations(t)
}
