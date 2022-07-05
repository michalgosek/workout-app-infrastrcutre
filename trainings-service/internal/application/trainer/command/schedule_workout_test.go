package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldScheduleWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	ctx := context.Background()
	workoutDate := time.Now().Add(time.Hour * 24)
	repository := new(mocks.TrainerRepository)
	SUT := command.NewScheduleWorkoutHandler(repository)

	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, mock.Anything).Return(nil)
	repository.EXPECT().QueryTrainerWorkoutGroupWithDate(ctx, trainerUUID, workoutDate).Return(trainer.WorkoutGroup{}, nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: trainerUUID,
		TrainerName: "John Doe",
		GroupName:   "dummy",
		GroupDesc:   "dummy",
		Date:        workoutDate,
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotScheduleWorkoutGroupWhenRepositoryFailureForUpsert_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	ctx := context.Background()
	workoutDate := time.Now().Add(time.Hour * 24)
	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	SUT := command.NewScheduleWorkoutHandler(repository)

	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, mock.Anything).Return(expectedError)
	repository.EXPECT().QueryTrainerWorkoutGroupWithDate(ctx, trainerUUID, workoutDate).Return(trainer.WorkoutGroup{}, nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: trainerUUID,
		TrainerName: "John Doe",
		GroupName:   "dummy",
		GroupDesc:   "dummy",
		Date:        workoutDate,
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldNotScheduleWorkoutGroupWhenRepositoryFailureForQuery_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	ctx := context.Background()
	workoutDate := time.Now().Add(time.Hour * 24)
	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	SUT := command.NewScheduleWorkoutHandler(repository)
	repository.EXPECT().QueryTrainerWorkoutGroupWithDate(ctx, trainerUUID, workoutDate).Return(trainer.WorkoutGroup{}, expectedError)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: trainerUUID,
		TrainerName: "John Doe",
		GroupName:   "dummy",
		GroupDesc:   "dummy",
		Date:        workoutDate,
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenAttemptsToScheduleWorkoutWithDuplicatedDate_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "401f8468-22eb-42bb-aacb-31774423d327"

	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	SUT := command.NewScheduleWorkoutHandler(repository)

	repository.EXPECT().QueryTrainerWorkoutGroupWithDate(ctx, workoutGroup.TrainerUUID(), workoutGroup.Date()).Return(workoutGroup, nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: workoutGroup.TrainerUUID(),
		TrainerName: workoutGroup.TrainerName(),
		GroupName:   workoutGroup.Name(),
		GroupDesc:   workoutGroup.Description(),
		Date:        workoutGroup.Date(),
	})

	// then:
	assertions.Equal(err, command.ErrWorkoutGroupDateDuplicated)
	repository.AssertExpectations(t)
}
