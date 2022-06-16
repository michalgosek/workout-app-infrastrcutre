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

func TestShouldGetRequestedTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	SUT := query.NewGetWorkoutHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, workout.UUID(), trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(workout, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldGetEmptyTrainerWorkoutGroupWhenRequestedGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := query.NewGetWorkoutHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenWhenRequestedGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const secondTrainerUUID = "b6236555-b4f5-4b2a-8c27-87f20ec71961"

	ctx := context.Background()
	secondTrainerGroup := testutil.NewTrainerWorkoutGroup(secondTrainerUUID)
	workoutUUID := secondTrainerGroup.UUID()
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(secondTrainerGroup, nil)
	SUT := query.NewGetWorkoutHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Equal(query.ErrWorkoutGroupNotOwner, err)
	assertions.Empty(actualSchedule)
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
	repository.EXPECT().QueryWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, expectedError)
	SUT := query.NewGetWorkoutHandler(repository)

	// when:
	_, err := SUT.Do(ctx, groupUUID, trainerUUID)

	// then:
	assertions.ErrorContains(err, err.Error())
	repository.AssertExpectations(t)
}
