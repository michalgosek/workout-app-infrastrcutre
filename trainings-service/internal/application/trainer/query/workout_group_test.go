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
	expectedGroup := query.WorkoutGroupDetails{
		TrainerUUID: group.TrainerUUID(),
		TrainerName: group.TrainerName(),
		GroupUUID:   group.UUID(),
		GroupDesc:   group.Description(),
		GroupName:   group.Name(),
		Customers:   query.ConvertToCustomersData(group.CustomerDetails()),
		Date:        group.Date().String(),
	}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, group.UUID()).Return(group, nil)
	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	actualGroup, err := SUT.Do(ctx, group.UUID(), trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
	repository.AssertExpectations(t)
}

func TestShouldGetEmptyTrainerWorkoutGroupWhenRequestedGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	group := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, workoutUUID).Return(group, nil)

	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	actualGroup, err := SUT.Do(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenWhenRequestedGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const secondTrainerUUID = "b6236555-b4f5-4b2a-8c27-87f20ec71961"

	ctx := context.Background()
	secondTrainerGroup := testutil.NewTrainerWorkoutGroup(secondTrainerUUID)
	groupUUID := secondTrainerGroup.UUID()

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, groupUUID).Return(secondTrainerGroup, nil)
	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	actualGroup, err := SUT.Do(ctx, groupUUID, trainerUUID)

	// then:
	assertions.Equal(query.ErrWorkoutGroupNotOwner, err)
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
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, expectedError)
	SUT := query.NewWorkoutGroupHandler(repository)

	// when:
	_, err := SUT.Do(ctx, groupUUID, trainerUUID)

	// then:
	assertions.ErrorContains(err, err.Error())
	repository.AssertExpectations(t)
}
