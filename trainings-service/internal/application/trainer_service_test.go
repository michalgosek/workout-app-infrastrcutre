package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
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

func TestShouldReturnErrorWhenTrainerSchedulesRepositoryFailure_Unit(t *testing.T) {
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

func TestShouldReturnTrainerScheduleWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	repository := new(mocks.TrainerRepository)
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	repository.EXPECT().QuerySchedule(ctx, expectedSchedule.UUID()).Return(expectedSchedule, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedule(ctx, expectedSchedule.UUID(), trainerUUID)

	// then:
	assert.Nil(err)
	assertTrainerSchedule(assert, expectedSchedule, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnEmptyScheduleWhenNotBelongingToTrainer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	otherTrainerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	otherTrainerSchedule := testutil.GenerateTrainerSchedule(otherTrainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QuerySchedule(ctx, otherTrainerSchedule.UUID()).Return(otherTrainerSchedule, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedule(ctx, otherTrainerSchedule.UUID(), trainerUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnEmptyTrainerSchedulesWhenNotExistWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	repository := new(mocks.TrainerRepository)
	expectedSchedules := []trainer.TrainerSchedule{}
	repository.EXPECT().QuerySchedules(ctx, trainerUUID).Return(expectedSchedules, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Empty(expectedSchedules, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnAllTrainerSchedulesWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	firstSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	secodSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	expectedSchedules := []trainer.TrainerSchedule{firstSchedule, secodSchedule}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QuerySchedules(ctx, trainerUUID).Return(expectedSchedules, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSchedules, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToScheduleNotBelongingToTrainer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	trainerUUID := "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	otherTrainerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	otherTrainerSchedule := testutil.GenerateTrainerSchedule(otherTrainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QuerySchedule(ctx, otherTrainerSchedule.UUID()).Return(otherTrainerSchedule, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, otherTrainerSchedule.UUID(), trainerUUID)

	// then:
	assert.ErrorIs(err, application.ErrScheduleNotOwner)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToSpecifiedScheduleWhenRepositoryFailure_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	trainerUUID := "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	trainerSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	trainerScheduledWithCustomer := trainerSchedule
	trainerScheduledWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QuerySchedule(ctx, trainerSchedule.UUID()).Return(trainerSchedule, nil)
	repository.EXPECT().UpsertSchedule(ctx, trainerScheduledWithCustomer).Return(expectedErr)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, trainerSchedule.UUID(), trainerUUID)

	// then:
	assert.ErrorIs(err, application.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldAssignCustomerToSpecifiedSchedule_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "094bb50a-7da3-461f-86f6-46d16c055e1e"
	trainerUUID := "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	trainerSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	trainerScheduledWithCustomer := trainerSchedule
	trainerScheduledWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QuerySchedule(ctx, trainerSchedule.UUID()).Return(trainerSchedule, nil)
	repository.EXPECT().UpsertSchedule(ctx, trainerScheduledWithCustomer).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, trainerSchedule.UUID(), trainerUUID)

	// then:
	assert.Nil(err)
	repository.AssertExpectations(t)
}

func assertTrainerSchedule(assert *assert.Assertions, expectedSchedule trainer.TrainerSchedule, actualSchedule trainer.TrainerSchedule) {
	assert.Equal(expectedSchedule.UUID(), actualSchedule.UUID())
	assert.Equal(expectedSchedule.TrainerUUID(), actualSchedule.TrainerUUID())
	assert.Equal(expectedSchedule.Customers(), actualSchedule.Customers())
	assert.Equal(expectedSchedule.Limit(), actualSchedule.Limit())
	assert.Equal(expectedSchedule.Desc(), actualSchedule.Desc())
}
