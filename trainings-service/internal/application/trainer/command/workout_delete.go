package command

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutDeleter interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	DeleteWorkoutGroup(ctx context.Context, groupUUID string) error
}

type WorkoutDeleteHandler struct {
	repository WorkoutDeleter
}

func (w *WorkoutDeleteHandler) Do(ctx context.Context, groupUUID, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutDeleteHandler"})
	group, err := w.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		const s = "query workout group UUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	if group.TrainerUUID() != trainerUUID {
		const s = "query workout group UUID: %s does not belong to trainerUUID: %s"
		logger.Errorf(s, groupUUID, trainerUUID)
		return ErrWorkoutGroupNotOwner
	}
	err = w.repository.DeleteWorkoutGroup(ctx, groupUUID)
	if err != nil {
		const s = "delete workout group UUID: %s for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, groupUUID, trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewWorkoutDeleteHandler(w WorkoutDeleter) *WorkoutDeleteHandler {
	if w == nil {
		panic("nil repository")
	}
	return &WorkoutDeleteHandler{
		repository: w,
	}
}
