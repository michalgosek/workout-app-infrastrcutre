package command_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUnassignCustomerHandler_ShouldUnassignCustomerFromWorkoutWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "728e2050-f5a4-4297-b613-2e506ce1555f"
		groupUUID    = "63fd9cb4-c4d1-4ebd-9335-465153f1aeb3"
		trainerUUID  = "0862326f-768f-4746-a6d5-e7c8b903be47"
	)
	ctx := context.Background()
	service := mocks.NewTrainingsService(t)
	SUT, _ := command.NewUnassignCustomerHandler(service)

	service.EXPECT().CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomerArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, service)
}

func TestUnassignCustomerHandler_ShouldNotUnassignCustomerFromWorkoutWhenTrainingsServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "728e2050-f5a4-4297-b613-2e506ce1555f"
		groupUUID    = "63fd9cb4-c4d1-4ebd-9335-465153f1aeb3"
		trainerUUID  = "0862326f-768f-4746-a6d5-e7c8b903be47"
	)
	ctx := context.Background()
	service := mocks.NewTrainingsService(t)
	SUT, _ := command.NewUnassignCustomerHandler(service)

	serviceFailureErr := errors.New("service failure err")
	service.EXPECT().CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	}).Return(serviceFailureErr)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomerArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, serviceFailureErr)
	mock.AssertExpectationsForObjects(t, service)
}
