package trainer_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	domain "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestService_ShouldCancelWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroup := newTestTrainerWorkoutGroup(trainerUUID)
	groupUUID := trainerWorkoutGroup.UUID()

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroup, nil)
	repository.EXPECT().DeleteTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(nil)

	// when:
	err := SUT.CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelWorkoutGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID   = "744b132b-837e-470d-9594-3be39832d87b"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupNotFound := domain.WorkoutGroup{}
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupNotFound, nil)

	// when:
	err := SUT.CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, trainer.ErrResourceNotFound)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelWorkoutGroupWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID   = "570a80c1-6f44-4154-b6fb-be099e461583"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(domain.WorkoutGroup{}, repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, trainer.ErrQueryTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelWorkoutGroupWhenDeleteTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroup := newTestTrainerWorkoutGroup(trainerUUID)
	groupUUID := trainerWorkoutGroup.UUID()

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroup, nil)
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().DeleteTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, trainer.ErrDeleteTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldAssignCustomerToWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupWithoutCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer := trainerWorkoutGroupWithoutCustomer
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	expectedWorkoutGroupDetails := trainer.AssignedCustomerWorkoutGroupDetails{
		UUID:        trainerWorkoutGroupWithCustomer.UUID(),
		TrainerUUID: trainerWorkoutGroupWithCustomer.TrainerUUID(),
		Name:        trainerWorkoutGroupWithCustomer.Name(),
		Date:        trainerWorkoutGroupWithCustomer.Date(),
	}

	groupUUID := trainerWorkoutGroupWithoutCustomer.UUID()
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithoutCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithCustomer).Return(nil)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
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

func TestService_ShouldNotAssignCustomerToWorkoutGroupWhenGroupNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		groupUUID    = "2212d1aa-ce01-4f32-bbf1-240ed66da5d3"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(domain.WorkoutGroup{}, nil)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrResourceNotFound)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotAssignCustomerToWorkoutWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
		groupUUID    = "de3ef438-6c06-43a8-b6fe-3f3a386c3054"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(domain.WorkoutGroup{}, repositoryFailureErr)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrQueryTrainerWorkoutGroup)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotAssignCustomerToWorkoutWhenUpsertTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupWithoutCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer := trainerWorkoutGroupWithoutCustomer
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	groupUUID := trainerWorkoutGroupWithoutCustomer.UUID()
	repositoryFailureErr := errors.New("repository failure")

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithoutCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithCustomer).Return(repositoryFailureErr)

	// when:
	actualWorkoutGroupDetails, err := SUT.AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrUpsertTrainerWorkoutGroup)
	assertions.Empty(actualWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldCancelCustomerWorkoutParticipationWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	trainerWorkoutGroupWithoutCustomer := trainerWorkoutGroupWithCustomer
	trainerWorkoutGroupWithoutCustomer.UnregisterCustomer(customerUUID)

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithoutCustomer).Return(nil)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelCustomerWorkoutParticipationWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(domain.WorkoutGroup{}, repositoryFailureErr)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrQueryTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelCustomerWorkoutParticipationWhenUpsertTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "c262fca7-1364-4a05-bf4a-bffd3c28698e"
		groupUUID    = "5c7b8314-9c03-49b6-9edf-ab4df4cb778b"
		customerUUID = "cf12e268-3733-499e-84c1-b7e19095fb03"
		customerName = "John Doe"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	trainerWorkoutGroupWithCustomer.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	trainerWorkoutGroupWithoutCustomer := trainerWorkoutGroupWithCustomer
	trainerWorkoutGroupWithoutCustomer.UnregisterCustomer(customerUUID)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainerWorkoutGroupWithCustomer, nil)
	repository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutGroupWithoutCustomer).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrUpsertTrainerWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldGetCustomerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "647909be-5eba-4ae1-9d33-dda8b734a9cc"
		customerUUID = "89629ed6-87e8-4f9e-a115-0c576a732b22"
		customerName = "John Doe"
	)
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	ctx := context.Background()
	expectedWorkoutGroup := newTestTrainerWorkoutGroup(trainerUUID)
	expectedWorkoutGroup.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))

	repository.EXPECT().QueryCustomerWorkoutGroup(ctx, trainerUUID, expectedWorkoutGroup.UUID(), customerUUID).Return(expectedWorkoutGroup, nil)

	// when:
	actualWorkoutGroup, err := SUT.GetCustomerWorkoutGroup(ctx, trainer.WorkoutGroupWithCustomerArgs{
		TrainerUUID:  trainerUUID,
		CustomerUUID: customerUUID,
		GroupUUID:    expectedWorkoutGroup.UUID(),
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkoutGroup, actualWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)

}

func TestService_ShouldNotGetCustomerWorkoutGroupWhenQueryCustomerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "647909be-5eba-4ae1-9d33-dda8b734a9cc"
		customerUUID = "89629ed6-87e8-4f9e-a115-0c576a732b22"
		groupUUID    = "ca6010c3-c3c2-486b-8d04-5ee5e497ca29"
	)
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutGroup(ctx, trainerUUID, groupUUID, customerUUID).Return(domain.WorkoutGroup{}, repositoryFailureErr)

	// when:
	actualWorkoutGroup, err := SUT.GetCustomerWorkoutGroup(ctx, trainer.WorkoutGroupWithCustomerArgs{
		TrainerUUID:  trainerUUID,
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Equal(err, trainer.ErrQueryTrainerWorkoutGroupWithCustomer)
	assertions.Empty(actualWorkoutGroup)
	mock.AssertExpectationsForObjects(t, repository)

}

func TestService_ShouldCancelWorkoutGroupsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "647909be-5eba-4ae1-9d33-dda8b734a9cc"
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	repository.EXPECT().DeleteTrainerWorkoutGroups(ctx, trainerUUID).Return(nil)

	// when:
	err := SUT.CancelWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestService_ShouldNotCancelWorkoutGroupsWithSuccessWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "647909be-5eba-4ae1-9d33-dda8b734a9cc"
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := trainer.NewTrainerService(repository)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().DeleteTrainerWorkoutGroups(ctx, trainerUUID).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.ErrorIs(err, trainer.ErrDeleteTrainerWorkoutGroups)
	mock.AssertExpectationsForObjects(t, repository)
}

func newTestTrainerWorkoutGroup(trainerUUID string) domain.WorkoutGroup {
	schedule := time.Now().AddDate(0, 0, 1)
	group, err := domain.NewWorkoutGroup(trainerUUID, "dummy_trainer", "dummy_group", "dummy_desc", schedule)
	if err != nil {
		panic(err)
	}
	return group
}

func newTestCustomerDetails(customerUUID, name string) customer.Details {
	details, err := customer.NewCustomerDetails(customerUUID, "John Doe")
	if err != nil {
		panic(err)
	}
	return details
}
