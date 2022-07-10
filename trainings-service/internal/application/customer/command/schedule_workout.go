package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
)

type ScheduleWorkoutArgs struct {
	CustomerUUID string
	CustomerName string
	TrainerUUID  string
	GroupUUID    string
}

type ScheduleWorkoutHandler struct {
	trainingsService TrainingsService
}

func (s *ScheduleWorkoutHandler) Do(ctx context.Context, args ScheduleWorkoutArgs) error {
	err := s.trainingsService.AssignCustomerToWorkoutGroup(ctx, trainings.AssignCustomerToWorkoutArgs{
		TrainerUUID:  args.TrainerUUID,
		GroupUUID:    args.GroupUUID,
		CustomerUUID: args.CustomerUUID,
		CustomerName: args.CustomerName,
	})
	if err != nil {
		return fmt.Errorf("trainings service  failure: %w", err)
	}
	return nil
}

func NewScheduleWorkoutHandler(t TrainingsService) (*ScheduleWorkoutHandler, error) {
	if t == nil {
		return nil, ErrNilTrainingsService
	}

	h := ScheduleWorkoutHandler{
		trainingsService: t,
	}
	return &h, nil
}
