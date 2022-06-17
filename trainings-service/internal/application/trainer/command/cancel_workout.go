package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkoutHandler struct {
	repository TrainerRepository
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, groupUUID, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutDeleteHandler"})
	group, err := c.repository.QueryTrainerWorkoutGroup(ctx, groupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	if group.TrainerUUID() != trainerUUID {
		logger.Errorf("workout group UUID: %s does not belong to trainerUUID: %s", groupUUID, trainerUUID)
		return ErrWorkoutGroupNotOwner
	}
	err = c.repository.DeleteTrainerWorkoutGroup(ctx, groupUUID)
	if err != nil {
		logger.Errorf("delete workout group UUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutHandler(w TrainerRepository) *CancelWorkoutHandler {
	if w == nil {
		panic("nil trainer repository")
	}
	return &CancelWorkoutHandler{
		repository: w,
	}
}
