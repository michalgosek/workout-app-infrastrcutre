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
	repository.EXPECT().UpsertTrainerSchedule(ctx, mock.Anything).Return(nil)
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
	repository.EXPECT().UpsertTrainerSchedule(ctx, mock.Anything).Return(expectedError)
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
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	repository.EXPECT().QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID()).Return(expectedSchedule, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedule(ctx, expectedSchedule.UUID(), trainerUUID)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnEmptyScheduleWhenNotBelongingToTrainer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const providedScheduleUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	emptySchedule := trainer.TrainerSchedule{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedule(ctx, providedScheduleUUID, trainerUUID).Return(emptySchedule, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetSchedule(ctx, providedScheduleUUID, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldReturnEmptyTrainerSchedulesWhenNotExistWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	expectedSchedules := []trainer.TrainerSchedule{}
	repository.EXPECT().QueryTrainerSchedules(ctx, trainerUUID).Return(expectedSchedules, nil)

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
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)
	expectedSchedules := []trainer.TrainerSchedule{first, second}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedules(ctx, trainerUUID).Return(expectedSchedules, nil)
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
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const providedScheduleUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	emptySchedule := trainer.TrainerSchedule{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedule(ctx, providedScheduleUUID, trainerUUID).Return(emptySchedule, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, providedScheduleUUID, trainerUUID)

	// then:
	assert.ErrorIs(err, application.ErrScheduleNotOwner)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToSpecifiedScheduleWhenRepositoryFailure_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	schedule := testutil.GenerateTrainerSchedule(trainerUUID)
	scheduledWithCustomer := schedule
	scheduledWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(schedule, nil)
	repository.EXPECT().UpsertTrainerSchedule(ctx, scheduledWithCustomer).Return(expectedErr)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, schedule.UUID(), trainerUUID)

	// then:
	assert.ErrorIs(err, application.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldAssignCustomerToSpecifiedSchedule_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	schedule := testutil.GenerateTrainerSchedule(trainerUUID)
	scheduledWithCustomer := schedule
	scheduledWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(schedule, nil)
	repository.EXPECT().UpsertTrainerSchedule(ctx, scheduledWithCustomer).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssingCustomer(ctx, customerUUID, schedule.UUID(), trainerUUID)

	// then:
	assert.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteScheduleNotOwnedByTrainer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const providedScheduleUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	emptySchedule := trainer.TrainerSchedule{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedule(ctx, providedScheduleUUID, trainerUUID).Return(emptySchedule, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteSchedule(ctx, providedScheduleUUID, trainerUUID)

	// then:
	assert.Equal(application.ErrScheduleNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteScheduleOwnedByTrainerWhenRepositoryFailure_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	schedule := testutil.GenerateTrainerSchedule(trainerUUID)
	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(schedule, nil)
	repository.EXPECT().CancelTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(expectedErr)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteSchedule(ctx, schedule.UUID(), schedule.TrainerUUID())

	// then:
	assert.Contains(err.Error(), expectedErr.Error())
	repository.AssertExpectations(t)
}

func TestShouldDeleteScheduleOwnedByTrainer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	schedule := testutil.GenerateTrainerSchedule(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(schedule, nil)
	repository.EXPECT().CancelTrainerSchedule(ctx, schedule.UUID(), schedule.TrainerUUID()).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteSchedule(ctx, schedule.UUID(), schedule.TrainerUUID())

	// then:
	assert.Nil(err)
	repository.AssertExpectations(t)
}
