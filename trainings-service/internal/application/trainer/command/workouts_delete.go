package command

import (
	"context"
	"fmt"
)

type WorkoutsDeleter interface {
	DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error
}

type WorkoutsDeleteHandler struct {
	repository WorkoutsDeleter
}

func (w *WorkoutsDeleteHandler) Do(ctx context.Context, trainerUUID string) error {
	err := w.repository.DeleteWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return fmt.Errorf("delete workout groups failed: %v", err)
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
