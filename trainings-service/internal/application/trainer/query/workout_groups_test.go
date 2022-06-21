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

func TestShouldGetEmptyTrainerWorkoutGroupsWhenNonOfGroupsDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	var groups []trainer.WorkoutGroup
	repository.EXPECT().QueryTrainerWorkoutGroups(ctx, trainerUUID).Return(groups, nil)

	SUT := query.NewWorkoutGroupsHandler(repository)

	// when:
	actualSchedule, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(groups, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldNotGetTrainerWorkoutGroupsWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroups(ctx, trainerUUID).Return(nil, expectedError)
	SUT := query.NewWorkoutGroupsHandler(repository)

	// when:
	_, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Equal(err, query.ErrRepositoryFailure)
	repository.AssertExpectations(t)

}

func TestShouldGetAllTrainerWorkoutGroupsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	firstWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	secondWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)

	groups := []trainer.WorkoutGroup{firstWorkout, secondWorkout}
	expectedWorkouts := query.WorkoutGroupsDetails{
		WorkoutGroups: []query.WorkoutGroupDetails{
			{
				TrainerUUID: firstWorkout.TrainerUUID(),
				TrainerName: firstWorkout.TrainerName(),
				GroupUUID:   firstWorkout.UUID(),
				GroupDesc:   firstWorkout.Description(),
				GroupName:   firstWorkout.Name(),
				Customers:   query.ConvertToCustomersData(firstWorkout.CustomerDetails()),
				Date:        firstWorkout.Date().String(),
			},
			{
				TrainerUUID: secondWorkout.TrainerUUID(),
				TrainerName: secondWorkout.TrainerName(),
				GroupUUID:   secondWorkout.UUID(),
				GroupDesc:   secondWorkout.Description(),
				GroupName:   secondWorkout.Name(),
				Customers:   query.ConvertToCustomersData(secondWorkout.CustomerDetails()),
				Date:        secondWorkout.Date().String(),
			},
		},
	}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerWorkoutGroups(ctx, trainerUUID).Return(groups, nil)
	SUT := query.NewWorkoutGroupsHandler(repository)

	// when:
	actualGroups, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkouts, actualGroups)
	repository.AssertExpectations(t)
}
