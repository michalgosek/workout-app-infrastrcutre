package customer

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"time"
)

type CancelWorkoutDayArgs struct {
	CustomerUUID string
	GroupUUID    string
}

type ScheduleWorkoutDayArgs struct {
	CustomerUUID string
	CustomerName string
	TrainerUUID  string
	GroupUUID    string
	Date         time.Time
}

func isCustomerWorkoutDayEmpty(customerWorkoutDay customer.WorkoutDay) bool {
	var empty customer.WorkoutDay
	return customerWorkoutDay == empty
}

//go:generate mockery --name=Repository --case underscore --with-expecter
type Repository interface {
	QueryCustomerWorkoutDay(ctx context.Context, customerUUID, GroupUUID string) (customer.WorkoutDay, error)
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, groupUUID string) error
	DeleteCustomersWorkoutDaysWithGroup(ctx context.Context, groupUUID string) error
	DeleteCustomersWorkoutDaysWithTrainer(ctx context.Context, trainerUUID string) error
}

type Service struct {
	repository Repository
}

func (s *Service) CancelWorkoutDaysWithGroup(ctx context.Context, groupUUID string) error {
	err := s.repository.DeleteCustomersWorkoutDaysWithGroup(ctx, groupUUID)
	if err != nil {
		return ErrDeleteCustomersWorkoutDaysWithGroup
	}
	return nil
}

func (s *Service) CancelWorkoutDaysWithTrainer(ctx context.Context, trainerUUID string) error {
	err := s.repository.DeleteCustomersWorkoutDaysWithTrainer(ctx, trainerUUID)
	if err != nil {
		return ErrDeleteCustomersWorkoutDaysWithTrainer
	}
	return nil
}

func (s *Service) CancelWorkoutDay(ctx context.Context, args CancelWorkoutDayArgs) error {
	customerWorkoutDay, err := s.repository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryCustomerWorkoutDay
	}
	if isCustomerWorkoutDayEmpty(customerWorkoutDay) {
		return nil
	}
	err = s.repository.DeleteCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrDeleteCustomerWorkoutDay
	}
	return nil
}

func (s *Service) ScheduleWorkoutDay(ctx context.Context, args ScheduleWorkoutDayArgs) error {
	customerWorkoutDay, err := s.repository.QueryCustomerWorkoutDay(ctx, args.CustomerUUID, args.GroupUUID)
	if err != nil {
		return ErrQueryCustomerWorkoutDay
	}
	if !isCustomerWorkoutDayEmpty(customerWorkoutDay) {
		return ErrResourceDuplicated
	}
	workoutDay, err := customer.NewWorkoutDay(args.CustomerUUID, args.CustomerName, args.GroupUUID, args.TrainerUUID, args.Date)
	if err != nil {
		return err
	}
	err = s.repository.UpsertCustomerWorkoutDay(ctx, workoutDay)
	if err != nil {
		return ErrUpsertCustomerWorkoutDay
	}
	return nil
}

func NewCustomerService(r Repository) (*Service, error) {
	if r == nil {
		return nil, ErrNilCustomerRepository
	}
	s := Service{repository: r}
	return &s, nil
}

var (
	ErrDeleteCustomersWorkoutDaysWithGroup   = errors.New("cmd delete customers workout days with group failure")
	ErrUpsertCustomerWorkoutDay              = errors.New("cmd upsert customer workout day failure")
	ErrDeleteCustomerWorkoutDay              = errors.New("cmd delete customer workout group failure")
	ErrDeleteCustomersWorkoutDaysWithTrainer = errors.New("cmd delete customers workout days with trainer failure")
	ErrQueryCustomerWorkoutDay               = errors.New("query customer workout day failure")
)

var (
	ErrNilCustomerRepository = errors.New("nil customer repository")
)

var (
	ErrResourceDuplicated = errors.New("resource duplicated")
	ErrResourceNotFound   = errors.New("resource not found")
)
