package query_test

import (
	"context"
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
