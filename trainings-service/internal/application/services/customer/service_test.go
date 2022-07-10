package customer_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/customer/mocks"
	domain "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
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

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	date := time.Now().AddDate(0, 0, 1)
	workoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)

	ctx := context.Background()
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(workoutDay, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(nil)

	// when:
	err := SUT.CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
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

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)

	ctx := context.Background()
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(domain.WorkoutDay{}, nil)

	// when:
	err := SUT.CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
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

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(domain.WorkoutDay{}, repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, customer.ErrQueryCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotCancelWorkoutDayWhenDeleteCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	date := time.Now().AddDate(0, 0, 1)
	workoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(workoutDay, nil)
	repository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.ErrorIs(err, customer.ErrDeleteCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldScheduleWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		customerName = "John Doe"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(domain.WorkoutDay{}, nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout domain.WorkoutDay) {
		assertions.Equal(groupUUID, workout.GroupUUID())
		assertions.Equal(customerUUID, workout.CustomerUUID())
		assertions.Equal(date, workout.Date())
	}).Return(nil)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
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

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(domain.WorkoutDay{}, repositoryFailureErr)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.ErrorIs(err, customer.ErrQueryCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldScheduleWorkoutDayWhenUpsertCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
		customerName = "John Doe"
	)

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(domain.WorkoutDay{}, nil)
	repository.EXPECT().UpsertCustomerWorkoutDay(ctx, mock.Anything).Run(func(ctx context.Context, workout domain.WorkoutDay) {
		assertions.Equal(groupUUID, workout.GroupUUID())
		assertions.Equal(customerUUID, workout.CustomerUUID())
		assertions.Equal(date, workout.Date())
	}).Return(repositoryFailureErr)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		CustomerName: customerName,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.ErrorIs(err, customer.ErrUpsertCustomerWorkoutDay)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotScheduleDuplicateWorkoutDay_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		customerUUID = "b1828953-c76f-4c29-85c5-065072d321d2"
		groupUUID    = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)

	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)
	ctx := context.Background()
	date := time.Now().AddDate(0, 0, 1)

	duplicateWorkoutDay := newTestCustomerWorkoutDay(customerUUID, groupUUID, date)
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, groupUUID).Return(duplicateWorkoutDay, nil)

	// when:
	err := SUT.ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUID,
		Date:         date,
	})

	// then:
	assertions.Equal(err, customer.ErrResourceDuplicated)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldCancelWorkoutDaysWithGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		groupUUID = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)

	repository.EXPECT().DeleteCustomersWorkoutDaysWithGroup(ctx, groupUUID).Return(nil)

	// when:
	err := SUT.CancelWorkoutDaysWithGroup(ctx, groupUUID)

	// then:
	assertions.Nil(err)
	mock.AssertExpectationsForObjects(t, repository)
}

func TestCustomerService_ShouldNotCancelWorkoutDaysWhenDeleteCustomersWorkoutDaysWithGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		groupUUID = "2c828785-8cd2-4e89-b198-55dac43f8317"
	)
	ctx := context.Background()
	repository := new(mocks.Repository)
	SUT, _ := customer.NewCustomerService(repository)

	repositoryFailureErr := errors.New("repository failure")
	repository.EXPECT().DeleteCustomersWorkoutDaysWithGroup(ctx, groupUUID).Return(repositoryFailureErr)

	// when:
	err := SUT.CancelWorkoutDaysWithGroup(ctx, groupUUID)

	// then:
	assertions.Equal(err, customer.ErrDeleteCustomersWorkoutDaysWithGroup)
	mock.AssertExpectationsForObjects(t, repository)
}

func newTestCustomerWorkoutDay(customerUUID, groupUUID string, date time.Time) domain.WorkoutDay {
	const name = "John Doe"
	workoutDay, err := domain.NewWorkoutDay(customerUUID, name, groupUUID, date)
	if err != nil {
		panic(err)
	}
	return workoutDay
}
