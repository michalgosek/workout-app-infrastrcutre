package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldCreateTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	args := application.TrainerSchedule{
		TrainerUUID: "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a",
		Name:        "dummy",
		Desc:        "dummy",
		Date:        time.Now().Add(time.Hour * 24),
	}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().UpsertWorkoutGroup(ctx, mock.Anything).Return(nil)
	SUT := application.NewTrainerService(repository)

	// when:
	_, err := SUT.CreateWorkoutGroup(ctx, args)

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotCreateTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

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
	repository.EXPECT().UpsertWorkoutGroup(ctx, mock.Anything).Return(expectedError)
	SUT := application.NewTrainerService(repository)

	// when:
	_, err := SUT.CreateWorkoutGroup(ctx, args)

	// then:
	assertions.ErrorContains(err, err.Error())
	repository.AssertExpectations(t)
}

func TestShouldGetRequestedTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetWorkoutGroup(ctx, workout.UUID(), trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(workout, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldGetEmptyTrainerWorkoutGroupWhenRequestedGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetWorkoutGroup(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldNotGetTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const groupUUID = "ef547a4e-f0ef-4282-a308-e985cbac2a01"
	const trainerUUID = "5a6bca90-a6d8-43d7-b1f8-069f9d5e846a"

	expectedError := errors.New("repository failure")
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, groupUUID).Return(trainer.WorkoutGroup{}, expectedError)
	SUT := application.NewTrainerService(repository)

	// when:
	_, err := SUT.GetWorkoutGroup(ctx, groupUUID, trainerUUID)

	// then:
	assertions.ErrorContains(err, err.Error())
	repository.AssertExpectations(t)
}

func TestShouldGetEmptyTrainerWorkoutGroupsWhenNonOfGroupsDoesNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	repository := new(mocks.TrainerRepository)
	var workouts []trainer.WorkoutGroup
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(workouts, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(workouts, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldGetAllTrainerWorkoutGroupsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	ctx := context.Background()

	first := testutil.NewTrainerWorkoutGroup(trainerUUID)
	second := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workouts := []trainer.WorkoutGroup{first, second}

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroups(ctx, trainerUUID).Return(workouts, nil)
	SUT := application.NewTrainerService(repository)

	// when:
	actualSchedule, err := SUT.GetWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(workouts, actualSchedule)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	ctx := context.Background()
	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssignCustomer(ctx, application.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workoutUUID,
	})

	// then:
	assertions.ErrorIs(err, application.ErrScheduleNotOwner)
	repository.AssertExpectations(t)
}

func TestShouldNotAssignCustomerToSpecifiedWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutdWithCustomer := workout
	workoutdWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutdWithCustomer).Return(expectedErr)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssignCustomer(ctx, application.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workout.UUID(),
	})

	// then:
	assertions.ErrorIs(err, application.ErrRepositoryFailure)
	repository.AssertExpectations(t)
}

func TestShouldAssignCustomerToSpecifiedWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"

	ctx := context.Background()
	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	workoutWithCustomer := workout
	_ = workoutWithCustomer.AssignCustomer(customerUUID)

	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().UpsertWorkoutGroup(ctx, workoutWithCustomer).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.AssignCustomer(ctx, application.WorkoutRegistration{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    workout.UUID(),
	})

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const workoutUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"

	workout := trainer.WorkoutGroup{}
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workoutUUID).Return(workout, nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteWorkoutGroup(ctx, workoutUUID, trainerUUID)

	// then:
	assertions.Equal(application.ErrScheduleNotOwner, err)
	repository.AssertExpectations(t)
}

func TestShouldNotDeleteTrainerWorkoutGroupWhenRepositoryFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	expectedErr := errors.New("repository failure")
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().DeleteWorkoutGroup(ctx, workout.UUID()).Return(expectedErr)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteWorkoutGroup(ctx, workout.UUID(), workout.TrainerUUID())

	// then:
	assertions.Contains(err.Error(), expectedErr.Error())
	repository.AssertExpectations(t)
}

func TestShouldDeleteWorkoutGroupOwnedByTrainerWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"

	workout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	repository := new(mocks.TrainerRepository)
	repository.EXPECT().QueryWorkoutGroup(ctx, workout.UUID()).Return(workout, nil)
	repository.EXPECT().DeleteWorkoutGroup(ctx, workout.UUID()).Return(nil)

	SUT := application.NewTrainerService(repository)

	// when:
	err := SUT.DeleteWorkoutGroup(ctx, workout.UUID(), workout.TrainerUUID())

	// then:
	assertions.Nil(err)
	repository.AssertExpectations(t)
}
