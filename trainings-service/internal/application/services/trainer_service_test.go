package services_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTrainerService_ShouldAssignCustomerToWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithoutCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer := trainerWorkoutGroupWithoutCustomer
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	expectedWorkoutGroupDetails := services.AssignedCustomerWorkoutGroupDetails{
		UUID:        trainerWorkoutGroupWithCustomer.UUID(),
		TrainerUUID: trainerWorkoutGroupWithCustomer.TrainerUUID(),
		Name:        trainerWorkoutGroupWithCustomer.Name(),
		Date:        trainerWorkoutGroupWithCustomer.Date(),
	}

	groupUUID := trainerWorkoutGroupWithoutCustomer.UUID()
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithoutCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithCustomer).Return(nil)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkoutGroupDetails, actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotAssignCustomerToWorkoutGroupWhenGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		groupUUID    = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, nil)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrResourceNotFound)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotAssignDuplicatedCustomerToWorkoutGroup_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	groupUUID := trainerWorkoutGroupWithCustomer.UUID()
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithCustomer, nil)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrResourceDuplicated)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotAssignCustomerToWorkoutWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
		groupUUID    = "de3ef438-6c06-43a8-b6fe-3f3a386c3054"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, repositoryFailureErr)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrQueryTrainerWorkoutGroup)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotAssignCustomerToWorkoutWhenUpsertTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithoutCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer := trainerWorkoutGroupWithoutCustomer
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	groupUUID := trainerWorkoutGroupWithoutCustomer.UUID()
	repositoryFailureErr := errors.New("repository failure")

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithoutCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithCustomer).Return(repositoryFailureErr)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrUpsertTrainerWorkoutGroup)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldCancelCustomerWorkoutParticipationWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	trainerWorkoutGroupWithoutCustomer := trainerWorkoutGroupWithCustomer
	trainerWorkoutGroupWithoutCustomer.UnregisterCustomer(customerUUID)

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithoutCustomer).Return(nil)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, services.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotCancelCustomerWorkoutParticipationWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, repositoryFailureErr)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, services.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrQueryTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestTrainerService_ShouldNotCancelCustomerWorkoutParticipationWhenUpsertTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.TrainerRepository)
	SUT, _ := services.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	trainerWorkoutGroupWithoutCustomer := trainerWorkoutGroupWithCustomer
	trainerWorkoutGroupWithoutCustomer.UnregisterCustomer(customerUUID)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithoutCustomer).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, services.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, services.ErrUpsertTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func newTestTrainerWorkoutGroup(trainerUUID string) trainer.WorkoutGroup {
	schedule := time.Now().AddDate(0, 0, 1)
	group, err := trainer.NewWorkoutGroup(trainerUUID, "dummy_trainer", "dummy_group", "dummy_desc", schedule)
	if err != nil {
		panic(err)
	}
	return group
}
