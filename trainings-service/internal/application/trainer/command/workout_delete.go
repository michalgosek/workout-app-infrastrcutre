package command

import (
	"context"
	"fmt"

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
	group, err := w.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("query trainer workout group failed: %v", err)
	}
	if group.TrainerUUID() != trainerUUID {
		return ErrScheduleNotOwner
	}
	err = w.repository.DeleteWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("delete workout group failed: %v", err)
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
