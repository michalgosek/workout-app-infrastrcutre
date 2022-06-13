package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldCreateTrainerScheduleWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	args := application.TrainerSchedule{
		TrainerUUID: "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a",
		Name:        "dummy",
		Desc:        "dummy",
		Date:        time.Now().Add(time.Hour * 24),
	}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().UpsertSchedule(ctx, mock.Anything).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.CreateTrainerSchedule(ctx, args)

	// then:
	assert.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenTrainerScheduleRepositoryFailure_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	args := application.TrainerSchedule{
		TrainerUUID: "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a",
		Name:        "dummy",
		Desc:        "dummy",
		Date:        time.Now().Add(time.Hour * 24),
	}

	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().UpsertSchedule(ctx, mock.Anything).Return(expectedError)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.CreateTrainerSchedule(ctx, args)

	// then:
	assert.ErrorContains(err, err.Error())
	repository.AssertExpectations(t)
}

func TestShouldReturnTrainerSessionWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	repository := new(mocks.TrainerRepository)
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	SUT := application.NewTrainerService(repository)

	repository.EXPECT().QuerySchedule(ctx, trainerUUID).Return(expectedSchedule, nil)

	// when:
	actualSchedule, err := SUT.GetSchedule(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
	repository.AssertExpectations(t)
}
