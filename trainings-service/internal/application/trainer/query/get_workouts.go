package query

import (
	"context"
	"fmt"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutGroupsGetter interface {
	QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
}

type GetWorkoutsHandler struct {
	repository WorkoutGroupsGetter
}

func (g *GetWorkoutsHandler) Do(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	groups, err := g.repository.QueryWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return nil, fmt.Errorf("get groups failed: %v", err)
	}
	return groups, nil
}

func NewGetWorkoutsHandler(w WorkoutGroupsGetter) *GetWorkoutsHandler {
	if w == nil {
		panic("nil repository")
	}
	return &GetWorkoutsHandler{
		repository: w,
	}
}
