package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
)

type ScheduleWorkout struct {
	CustomerUUID string
	CustomerName string
	TrainerUUID  string
	GroupUUID    string
}

type ScheduleWorkoutHandler struct {
	customerService CustomerService
	trainerService  TrainerService
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, w ScheduleWorkout) error {
	details, err := s.trainerService.AssignCustomerToWorkoutGroup(ctx, services.AssignCustomerToWorkoutGroupArgs{
		TrainerUUID:  w.TrainerUUID,
		GroupUUID:    w.GroupUUID,
		CustomerUUID: w.CustomerUUID,
		CustomerName: w.CustomerName,
	})
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	err = s.customerService.ScheduleWorkoutDay(ctx, services.ScheduleWorkoutDayArgs{
		CustomerUUID: w.CustomerUUID,
		GroupUUID:    w.GroupUUID,
		Date:         details.Date,
	})
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	return nil
}

func NewScheduleWorkoutHandler(t TrainerService, c CustomerService) (*ScheduleWorkoutHandler, error) {
	if t == nil {
		return nil, ErrNilTrainerService
	}
	if c == nil {
		return nil, ErrNiLCustomerService
	}
	h := ScheduleWorkoutHandler{
		trainerService:  t,
		customerService: c,
	}
	return &h, nil
}

var (
	ErrNilTrainerService  = errors.New("nil trainer service")
	ErrNiLCustomerService = errors.New("nil customer service")
)
