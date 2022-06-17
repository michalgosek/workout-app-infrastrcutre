package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type CancelWorkoutHandlerRepository interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	DeleteWorkoutGroup(ctx context.Context, groupUUID string) error
}

type CancelWorkoutHandler struct {
	repository CancelWorkoutHandlerRepository
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, groupUUID, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutDeleteHandler"})
	group, err := c.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		const s = "query workout group UUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	if group.TrainerUUID() != trainerUUID {
		const s = "workout group UUID: %s does not belong to trainerUUID: %s"
		logger.Errorf(s, groupUUID, trainerUUID)
		return ErrWorkoutGroupNotOwner
	}
	err = c.repository.DeleteWorkoutGroup(ctx, groupUUID)
	if err != nil {
		const s = "delete workout group UUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewCancelWorkoutHandler(w CancelWorkoutHandlerRepository) *CancelWorkoutHandler {
	if w == nil {
		panic("nil repository")
	}
	return &CancelWorkoutHandler{
		repository: w,
	}
}
