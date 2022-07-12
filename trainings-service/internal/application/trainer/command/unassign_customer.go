package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	"github.com/sirupsen/logrus"
)

type UnassignCustomerArgs struct {
	CustomerUUID string
	GroupUUID    string
	TrainerUUID  string
}

type UnassignCustomerHandler struct {
	trainingsService TrainingsService
}

func (u *UnassignCustomerHandler) Do(ctx context.Context, args UnassignCustomerArgs) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "UnregisterCustomerHandler"})
	err := u.trainingsService.CancelCustomerWorkout(ctx, trainings.CancelCustomerWorkoutArgs{
		CustomerUUID: args.CustomerUUID,
		GroupUUID:    args.GroupUUID,
		TrainerUUID:  args.TrainerUUID,
	})
	if err != nil {
		logger.Errorf("CMD - cancel customer workout failure: %s", err)
		return fmt.Errorf("trainings service failure: %w", err)
	}
	return nil
}

func NewUnassignCustomerHandler(t TrainingsService) (*UnassignCustomerHandler, error) {
	if t == nil {
		return &UnassignCustomerHandler{}, ErrNilTrainingsService
	}
	h := UnassignCustomerHandler{trainingsService: t}
	return &h, nil
}
