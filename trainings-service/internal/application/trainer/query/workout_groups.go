package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupsHandlerRepository interface {
	QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
}

type WorkoutGroupsHandler struct {
	repository WorkoutGroupsHandlerRepository
}

func (w *WorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Component": "GetWorkoutsHandler"})
	groups, err := w.repository.QueryWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		const s = "query workout groups for trainerUUID: %s failed: %v"
		logger.Errorf(s, trainerUUID, err)
		return nil, ErrRepositoryFailure
	}
	return groups, nil
}

func NewWorkoutGroupsHandler(w WorkoutGroupsHandlerRepository) *WorkoutGroupsHandler {
	if w == nil {
		panic("nil repository")
	}
	return &WorkoutGroupsHandler{
		repository: w,
	}
}
