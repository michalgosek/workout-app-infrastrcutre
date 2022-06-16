package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldGetEmptyTrainerWorkoutGroupsWhenNonOfGroupsDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	var workouts []trainer.WorkoutGroup
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(workouts, nil)

	SUT := query.NewGetWorkoutsHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(workouts, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldNotGetTrainerWorkoutGroupsWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(nil, expectedError)
	SUT := query.NewGetWorkoutsHandler(repository)

	// when:
	_, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.ErrorIs(err, query.ErrRepositoryFailure)
	repository.AssertExpectations(t)

}

func TestShouldGetAllTrainerWorkoutGroupsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	first := testutil.NewTrainerWorkoutGroup(trainerUUID)
	second := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workouts := []trainer.WorkoutGroup{first, second}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(workouts, nil)
	SUT := query.NewGetWorkoutsHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(workouts, actualSchedule)
	repository.AssertExpectations(t)
}
