package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkoutsHandlerRepository interface {
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type CancelWorkoutsHandler struct {
	repository CancelWorkoutsHandlerRepository
}

func (c *CancelWorkoutsHandler) Do(ctx context.Context, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutsDeleteHandler"})
	err := c.repository.DeleteWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		const s = "delete workout groups for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutsHandler(w CancelWorkoutsHandlerRepository) *CancelWorkoutsHandler {
	if w == nil {
		panic("nil repository")
	}
	return &CancelWorkoutsHandler{
		repository: w,
	}
}
