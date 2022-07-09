package services

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"time"
)

//go:generate mockery --name=CustomerRepository --case underscore --with-expecter
type CustomerRepository interface {
	QueryCustomerWorkoutDay(ctx context.Context, customerUUID, GroupUUID string) (customer.WorkoutDay, error)
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, groupUUID string) error
}

type CustomerService struct {
	repository CustomerRepository
}

type CancelWorkoutDayArgs struct {
	CustomerUUID string
	GroupUUID    string
}

func isCustomerWorkoutDayEmpty(customerWorkoutDay customer.WorkoutDay) bool {
	var empty customer.WorkoutDay
	return customerWorkoutDay == empty
}

func (c *CustomerService) CancelWorkoutDay(ctx context.Context, args CancelWorkoutDayArgs) error {
	customerWorkoutDay, err := c.repository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryCustomerWorkoutDay
	}
	if isCustomerWorkoutDayEmpty(customerWorkoutDay) {
		return nil
	}
	err = c.repository.DeleteCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrDeleteCustomerWorkoutDay
	}
	return nil
}

type ScheduleWorkoutDayArgs struct {
	CustomerUUID string
	GroupUUID    string
	Date         time.Time
}

func (c *CustomerService) ScheduleWorkoutDay(ctx context.Context, args ScheduleWorkoutDayArgs) error {
	customerWorkoutDay, err := c.repository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryCustomerWorkoutDay
	}
	if !isCustomerWorkoutDayEmpty(customerWorkoutDay) {
		return ErrResourceDuplicated
	}

	workoutDay, err := customer.NewWorkoutDay(args.CustomerUUID, args.GroupUUID, args.Date)
	if err != nil {
		return err
	}
	err = c.repository.UpsertCustomerWorkoutDay(ctx, workoutDay)
	if err != nil {
		return ErrUpsertCustomerWorkoutDay
	}
	return nil
}

func NewCustomerService(c CustomerRepository) (*CustomerService, error) {
	if c == nil {
		return nil, ErrNilCustomerRepository
	}
	s := CustomerService{repository: c}
	return &s, nil
}

var (
	ErrUpsertCustomerWorkoutDay = errors.New("cmd upsert customer workout day failure")
	ErrDeleteCustomerWorkoutDay = errors.New("cmd delete customer workout group failure")
	ErrQueryCustomerWorkoutDay  = errors.New("query customer workout day failure")
)

var (
	ErrNilCustomerRepository = errors.New("nil customer repository")
)

var (
	ErrResourceDuplicated = errors.New("resource duplicated")
	ErrResourceNotFound   = errors.New("resource not found")
)
