package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	"github.com/sirupsen/logrus"
)

type CancelWorkoutArgs struct {
	GroupUUID   string
	TrainerUUID string
}

type CancelWorkoutHandler struct {
	trainingsService TrainingsService
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkoutArgs) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutDeleteHandler"})
	err := c.trainingsService.CancelTrainerWorkoutGroup(ctx, trainings.CancelTrainerWorkoutGroupArgs{
		TrainerUUID: w.TrainerUUID,
		GroupUUID:   w.GroupUUID,
	})
	if err != nil {
		logger.Errorf("CMD - cancel trainer workout group Failure: %s", err)
		return fmt.Errorf("trainings service failure: %w", err)
	}
	return nil
}

func NewCancelWorkoutHandler(t TrainingsService) (*CancelWorkoutHandler, error) {
	if t == nil {
		return nil, ErrNilTrainingsService
	}
	h := CancelWorkoutHandler{
		trainingsService: t,
	}
	return &h, nil
}

var ErrNilTrainingsService = errors.New("nil trainings service")
