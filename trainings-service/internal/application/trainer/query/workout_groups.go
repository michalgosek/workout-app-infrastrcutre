package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupsHandler struct {
	repository TrainerRepository
}

func (w *WorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "GetWorkoutsHandler"})
	groups, err := w.repository.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("query workout groups for trainerUUID: %s failed: %v", trainerUUID, err)
		return nil, ErrRepositoryFailure
	}
	return groups, nil
}

func NewWorkoutGroupsHandler(t TrainerRepository) *WorkoutGroupsHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &WorkoutGroupsHandler{
		repository: t,
	}
}
