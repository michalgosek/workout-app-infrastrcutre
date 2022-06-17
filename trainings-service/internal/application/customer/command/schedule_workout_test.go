package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldScheduleWorkoutToSpecifiedWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.ScheduleWorkoutHandlerRepository)
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup
	workoutGroupWithCustomer.AssignCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutGroupWithCustomer).Return(nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout customer.WorkoutDay) {
		assertions.Equal(workoutGroup.UUID(), workout.GroupUUID())
		assertions.Equal(workoutGroup.Date(), workout.Date())
		assertions.Equal(customerUUID, workout.CustomerUUID())
	}).Return(nil)

	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:

	err := SUT.Do(ctx, command.WorkoutScheduleDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenQueryWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.ScheduleWorkoutHandlerRepository)
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const groupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	repository.EXPECT().QueryWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:

	err := SUT.Do(ctx, command.WorkoutScheduleDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.ScheduleWorkoutHandlerRepository)
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup
	workoutGroupWithCustomer.AssignCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutGroupWithCustomer).Return(errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:

	err := SUT.Do(ctx, command.WorkoutScheduleDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenUpsertCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.ScheduleWorkoutHandlerRepository)
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	workoutGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutGroupWithCustomer := workoutGroup
	workoutGroupWithCustomer.AssignCustomer(customerUUID)

	repository.EXPECT().QueryWorkoutGroup(ctx, workoutGroup.UUID()).Return(workoutGroup, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutGroupWithCustomer).Return(nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout customer.WorkoutDay) {
		assertions.Equal(workoutGroup.UUID(), workout.GroupUUID())
		assertions.Equal(workoutGroup.Date(), workout.Date())
		assertions.Equal(customerUUID, workout.CustomerUUID())
	}).Return(errors.New("err"))

	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:

	err := SUT.Do(ctx, command.WorkoutScheduleDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    workoutGroup.UUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenWorkoutGroupDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	repository := new(mocks.ScheduleWorkoutHandlerRepository)
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const workoutGroupUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	repository.EXPECT().QueryWorkoutGroup(ctx, workoutGroupUUID).Return(trainer.WorkoutGroup{}, nil)
	SUT := command.NewScheduleWorkoutHandler(repository)

	// when:
	err := SUT.Do(ctx, command.WorkoutScheduleDetails{
		CustomerUUID: customerUUID,
		GroupUUID:    workoutGroupUUID,
	})

	// then:
	assertions.Equal(err, command.ErrResourceNotFound)
	repository.AssertExpectations(t)
}
