package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type WorkoutsDeleter interface {
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type WorkoutsDeleteHandler struct {
	repository WorkoutsDeleter
}

func (w *WorkoutsDeleteHandler) Do(ctx context.Context, trainerUUID string) error {
	logger := logrus.WithFields(logrus.Fields{"Component": "WorkoutsDeleteHandler"})
	err := w.repository.DeleteWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		const s = "delete workout groups for trainerUUID: %s failed, reason: %v"
		logger.Errorf(s, trainerUUID, err)
		return ErrRepositoryFailure
	}
	return nil
}

func NewWorkoutsDeleteHandler(w WorkoutsDeleter) *WorkoutsDeleteHandler {
	if w == nil {
		panic("nil repository")
	}
	return &WorkoutsDeleteHandler{
		repository: w,
	}
}
