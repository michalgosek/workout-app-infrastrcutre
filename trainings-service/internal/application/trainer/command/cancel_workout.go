package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CancelWorkoutHandler struct {
	repository TrainerRepository
}

type CancelWorkout struct {
	GroupUUID   string
	TrainerUUID string
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkout) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutDeleteHandler"})
	group, err := c.repository.QueryTrainerWorkoutGroup(ctx, w.GroupUUID)
	if err != nil {
		logger.Errorf("query workout group UUID: %s for trainerUUID: %s failed, reason: %v", w.GroupUUID, w.TrainerUUID, err)
		return ErrRepositoryFailure
	}
	if group.TrainerUUID() != w.TrainerUUID {
		logger.Errorf("workout group UUID: %s does not belong to trainerUUID: %s", w.GroupUUID, w.TrainerUUID)
		return ErrWorkoutGroupNotOwner
	}
	err = c.repository.DeleteTrainerWorkoutGroup(ctx, w.GroupUUID)
	if err != nil {
		logger.Errorf("delete workout group UUID: %s for trainerUUID: %s failed, reason: %v", w.GroupUUID, w.TrainerUUID, err)
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
