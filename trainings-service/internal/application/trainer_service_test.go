package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldCreateTrainerScheduleWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	args := application.TrainerScheduleArgs{
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
