package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
)

type CancelWorkoutHandler struct {
	customerService CustomerService
	trainerService  TrainerService
}

type CancelWorkout struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkout) error {
	err := c.customerService.CancelWorkoutDay(ctx, services.CancelWorkoutDayArgs{
		CustomerUUID: w.CustomerUUID,
		GroupUUID:    w.GroupUUID,
	})
	if err != nil {
		return fmt.Errorf("customer service failure: %w", err)
	}
	err = c.trainerService.CancelCustomerWorkoutParticipation(ctx, services.CancelCustomerWorkoutParticipationArgs{
		CustomerUUID: w.CustomerUUID,
		GroupUUID:    w.GroupUUID,
		TrainerUUID:  w.TrainerUUID,
	})
	if err != nil {
		return fmt.Errorf("trainer service failure: %w", err)
	}
	return nil
}

func NewCancelWorkoutHandler(c CustomerService, t TrainerService) *CancelWorkoutHandler {
	if c == nil {
		panic("nil customer service") //fixme: should be returned err instead of panic!
	}
	if t == nil {
		panic("nil trainer service") //fixme: should be returned err instead of panic!
	}
	h := CancelWorkoutHandler{customerService: c, trainerService: t}
	return &h
}
