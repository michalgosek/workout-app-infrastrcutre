package command_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestScheduleWorkoutHandler_ShouldScheduleTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "d87493a3-8cb2-4701-8cd8-e13518c1d3f5"
		trainerName = "John Doe"
		groupName   = "group name"
		groupDesc   = "group desc"
	)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)
	service := mocks.NewTrainerService(t)
	SUT, _ := command.NewScheduleWorkoutHandler(service)

	service.EXPECT().CreateWorkoutGroup(ctx, trainer.CreateWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		TrainerName: trainerName,
		GroupName:   groupName,
		GroupDesc:   groupDesc,
		Date:        date,
	}).Return(nil)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkoutArgs{
		TrainerUUID: trainerUUID,
		TrainerName: trainerName,
		GroupName:   groupName,
		GroupDesc:   groupDesc,
		Date:        date,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, service)
}

func TestScheduleWorkoutHandler_ShouldNotScheduleTrainerWorkoutGroupWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "d87493a3-8cb2-4701-8cd8-e13518c1d3f5"
		trainerName = "John Doe"
		groupName   = "group name"
		groupDesc   = "group desc"
	)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)
	service := mocks.NewTrainerService(t)
	SUT, _ := command.NewScheduleWorkoutHandler(service)

	repositoryFailureErr := errors.New("repository failure error")
	service.EXPECT().CreateWorkoutGroup(ctx, trainer.CreateWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		TrainerName: trainerName,
		GroupName:   groupName,
		GroupDesc:   groupDesc,
		Date:        date,
	}).Return(repositoryFailureErr)

	// when:
	err := SUT.Do(ctx, command.ScheduleWorkoutArgs{
		TrainerUUID: trainerUUID,
		TrainerName: trainerName,
		GroupName:   groupName,
		GroupDesc:   groupDesc,
		Date:        date,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	mock.AssertExpectationsForObjects(t, service)
}
