package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkoutsHandler struct {
	repository TrainerRepository
}

func (c *CancelWorkoutsHandler) Do(ctx context.Context, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Trainer-CMD": "WorkoutsDeleteHandler"})
	err := c.repository.DeleteTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("delete workout groups for trainerUUID: %s failed, reason: %v", trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutsHandler(t TrainerRepository) *CancelWorkoutsHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &CancelWorkoutsHandler{
		repository: t,
	}
}
