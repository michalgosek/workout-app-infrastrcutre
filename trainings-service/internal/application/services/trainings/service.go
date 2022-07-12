package trainings

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
)

type CancelCustomerWorkoutArgs struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}
type AssignCustomerToWorkoutArgs struct {
	TrainerUUID  string
	GroupUUID    string
	CustomerUUID string
	CustomerName string
}

type CancelTrainerWorkoutGroupArgs struct {
	TrainerUUID string
	GroupUUID   string
}

//go:generate mockery --name=CustomerService --case underscore --with-expecter
type CustomerService interface {
	CancelWorkoutDay(ctx context.Context, args customer.CancelWorkoutDayArgs) error
	ScheduleWorkoutDay(ctx context.Context, args customer.ScheduleWorkoutDayArgs) error
	CancelWorkoutDaysWithGroup(ctx context.Context, groupUUID string) error
	CancelWorkoutDaysWithTrainer(ctx context.Context, trainerUUID string) error
}

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	AssignCustomerToWorkoutGroup(ctx context.Context, args trainer.AssignCustomerToWorkoutGroupArgs) (trainer.AssignedCustomerWorkoutGroupDetails, error)
	CancelCustomerWorkoutParticipation(ctx context.Context, args trainer.CancelCustomerWorkoutParticipationArgs) error
	CancelWorkoutGroup(ctx context.Context, args trainer.CancelWorkoutGroupArgs) error
	CancelWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type Service struct {
	customerService CustomerService
	trainerService  TrainerService
}

func (s *Service) CancelTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error {
	err := s.trainerService.CancelWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	err = s.customerService.CancelWorkoutDaysWithTrainer(ctx, trainerUUID)
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	return nil
}

func (s *Service) CancelTrainerWorkoutGroup(ctx context.Context, args CancelTrainerWorkoutGroupArgs) error {
	err := s.trainerService.CancelWorkoutGroup(ctx, trainer.CancelWorkoutGroupArgs{
		TrainerUUID: args.TrainerUUID,
		GroupUUID:   args.GroupUUID,
	})
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	err = s.customerService.CancelWorkoutDaysWithGroup(ctx, args.GroupUUID)
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	return nil
}

func (s *Service) CancelCustomerWorkout(ctx context.Context, args CancelCustomerWorkoutArgs) error {
	err := s.customerService.CancelWorkoutDay(ctx, customer.CancelWorkoutDayArgs{
		CustomerUUID: args.CustomerUUID,
		GroupUUID:    args.GroupUUID,
	})
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	err = s.trainerService.CancelCustomerWorkoutParticipation(ctx, trainer.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: args.CustomerUUID,
		GroupUUID:    args.GroupUUID,
		TrainerUUID:  args.TrainerUUID,
	})
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	return nil
}

func (s *Service) AssignCustomerToWorkoutGroup(ctx context.Context, args AssignCustomerToWorkoutArgs) error {
	details, err := s.trainerService.AssignCustomerToWorkoutGroup(ctx, trainer.AssignCustomerToWorkoutGroupArgs{
		TrainerUUID:  args.TrainerUUID,
		GroupUUID:    args.GroupUUID,
		CustomerUUID: args.CustomerUUID,
		CustomerName: args.CustomerName,
	})
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	err = s.customerService.ScheduleWorkoutDay(ctx, customer.ScheduleWorkoutDayArgs{
		CustomerUUID: args.CustomerUUID,
		CustomerName: args.CustomerName,
		GroupUUID:    args.GroupUUID,
		TrainerUUID:  args.TrainerUUID,
		Date:         details.Date,
	})
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	return nil
}

func NewService(c CustomerService, t TrainerService) (*Service, error) {
	if c == nil {
		return nil, ErrNilCustomerService
	}
	if t == nil {
		return nil, ErrNilTrainerService
	}
	s := Service{trainerService: t, customerService: c}
	return &s, nil
}

var (
	ErrNilTrainerService  = errors.New("nil trainer service")
	ErrNilCustomerService = errors.New("nil customer service")
)
