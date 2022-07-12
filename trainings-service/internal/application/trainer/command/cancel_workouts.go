package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

type CancelWorkoutsHandler struct {
	trainingsService TrainingsService
}

func (c *CancelWorkoutsHandler) Do(ctx context.Context, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "WorkoutsDeleteHandler"})
	err := c.trainingsService.CancelTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("CMD - cancel workout groups failure: %s", err)
		return fmt.Errorf("trainer service failure: %w", err)
	}
	return nil
}

func NewCancelWorkoutsHandler(t TrainingsService) (*CancelWorkoutsHandler, error) {
	if t == nil {
		return nil, ErrNilTrainingsService
	}
	h := CancelWorkoutsHandler{
		trainingsService: t,
	}
	return &h, nil
}

var ErrNilTrainerService = errors.New("nil trainer service")
