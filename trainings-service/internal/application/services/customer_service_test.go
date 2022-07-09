package services_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCustomerService_ShouldCancelWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	date := time.Now().AddDate(0, 0, 1)
	workoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)

	ctx := context.Background()
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(workoutDay, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(nil)

	// when:
	err := SUT.CancelWorkoutDay(ctx, services.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotCancelNonExistingWorkoutDay_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)

	ctx := context.Background()
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, nil)

	// when:
	err := SUT.CancelWorkoutDay(ctx, services.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotCancelWorkoutDayWhenQueryCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutDay(ctx, services.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, services.ErrQueryCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotCancelWorkoutDayWhenDeleteCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	date := time.Now().AddDate(0, 0, 1)
	workoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(workoutDay, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutDay(ctx, services.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, services.ErrDeleteCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldScheduleWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout customer.WorkoutDay) {
		assertions.Equal(groupUUID, workout.GroupUUID())
		assertions.Equal(customerUUID, workout.CustomerUUID())
		assertions.Equal(date, workout.Date())
	}).Return(nil)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotScheduleWorkoutDayWhenQueryCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, repositoryFailureErr)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.ErrorIs(err, services.ErrQueryCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldScheduleWorkoutDayWhenUpsertCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(customer.WorkoutDay{}, nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout customer.WorkoutDay) {
		assertions.Equal(groupUUID, workout.GroupUUID())
		assertions.Equal(customerUUID, workout.CustomerUUID())
		assertions.Equal(date, workout.Date())
	}).Return(repositoryFailureErr)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.ErrorIs(err, services.ErrUpsertCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotScheduleDuplicateWorkoutDay_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.CustomerRepository)
	SUT, _ := services.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	duplicateWorkoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(duplicateWorkoutDay, nil)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.Equal(err, services.ErrResourceDuplicated)
	mock.AssertExpectationsForObjects(t, repository)
}

func newTestCustomerWorkoutDay(customerUUID, groupUUID string, date time.Time) customer.WorkoutDay {
	workoutDay, err := customer.NewWorkoutDay(customerUUID, groupUUID, date)
	if err != nil {
		panic(err)
	}
	return workoutDay
}

func newTestCustomerDetails(customerUUID, name string) customer.Details {
	details, err := customer.NewCustomerDetails(customerUUID, "John Doe")
	if err != nil {
		panic(err)
	}
	return details
}
