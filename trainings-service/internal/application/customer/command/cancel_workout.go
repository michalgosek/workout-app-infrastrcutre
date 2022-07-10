package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
)

type CancelWorkoutHandler struct {
	trainingsService TrainingsService
}

type CancelWorkoutArgs struct {
	CustomerUUID string
	TrainerUUID  string
	GroupUUID    string
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, args CancelWorkoutArgs) error {
	err := c.trainingsService.CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: args.CustomerUUID,
		GroupUUID:    args.GroupUUID,
		TrainerUUID:  args.TrainerUUID,
	})
	if err != nil {
		return fmt.Errorf("trainings service failure: %w", err)
	}
	return nil
}

func NewCancelWorkoutHandler(t TrainingsService) (*CancelWorkoutHandler, error) {
	if t == nil {
		return &CancelWorkoutHandler{}, ErrNilTrainingsService
	}
	h := CancelWorkoutHandler{trainingsService: t}
	return &h, nil
}
