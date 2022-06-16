package query

import (
	"context"
	"fmt"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutGetter interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
}

type GetWorkoutHandler struct {
	repository WorkoutGetter
}

func (t *GetWorkoutHandler) Do(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	group, err := t.repository.QueryWorkoutGroup(ctx, groupUUID)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("query workout group failed: %v", err)
	}
	if group.TrainerUUID() != trainerUUID {
		return trainer.WorkoutGroup{}, nil
	}
	return group, nil
}

func NewGetWorkoutHandler(w WorkoutGetter) *GetWorkoutHandler {
	if w == nil {
		panic("nil repository")
	}
	return &GetWorkoutHandler{
		repository: w,
	}
}
