package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupsGetter interface {
	QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
}

type GetWorkoutsHandler struct {
	repository WorkoutGroupsGetter
}

func (g *GetWorkoutsHandler) Do(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "GetWorkoutsHandler"})
	groups, err := g.repository.QueryWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		const s = "query workout groups for trainerUUID: %s failed: %v"
		logger.Errorf(s, trainerUUID, err)
		return nil, ErrRepositoryFailure
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
