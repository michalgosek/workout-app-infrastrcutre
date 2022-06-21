package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldScheduleWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, mock.Anything).Return(nil)
	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a",
		TrainerName: "John Doe",
		GroupName:   "dummy",
		GroupDesc:   "dummy",
		Date:        time.Now().Add(time.Hour * 24),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotScheduleWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, mock.Anything).Return(expectedError)
	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkout{
		TrainerUUID: "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a",
		TrainerName: "John Doe",
		GroupName:   "dummy",
		GroupDesc:   "dummy",
		Date:        time.Now().Add(time.Hour * 24),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}
