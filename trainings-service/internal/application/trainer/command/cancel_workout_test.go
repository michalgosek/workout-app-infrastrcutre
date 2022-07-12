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

func TestCancelWorkoutHandler_ShouldCancelTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		groupUUID   = "33c12629-a38f-4437-8de1-6b1cb1b54fd1"
		trainerUUID = "b8c77ceb-158f-497e-8d2f-6d1764c7a607"
	)
	ctx := context.Background()
	service := mocks.NewTrainingsService(t)
	SUT, _ := command.NewCancelWorkoutHandler(service)

	service.EXPECT().CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutArgs{
		GroupUUID:   groupUUID,
		TrainerUUID: trainerUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, service)
}

func TestCancelWorkoutHandler_ShouldNotCancelTrainerWorkoutGroupWhenTrainingsServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		groupUUID   = "33c12629-a38f-4437-8de1-6b1cb1b54fd1"
		trainerUUID = "b8c77ceb-158f-497e-8d2f-6d1764c7a607"
	)
	ctx := context.Background()
	service := mocks.NewTrainingsService(t)
	SUT, _ := command.NewCancelWorkoutHandler(service)

	serviceFailureErr := errors.New("service failure")
	service.EXPECT().CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	}).Return(serviceFailureErr)

	// when:
	err := SUT.Do(ctx, command.CancelWorkoutArgs{
		GroupUUID:   groupUUID,
		TrainerUUID: trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, serviceFailureErr)
	mock.AssertExpectationsForObjects(t, service)
}
