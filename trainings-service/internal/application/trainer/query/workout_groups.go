package query

import (
	"context"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupsDetails struct {
	WorkoutGroups []WorkoutGroupDetails
}

type WorkoutGroupsHandler struct {
	repository TrainerRepository
}

func (w *WorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) (WorkoutGroupsDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "GetWorkoutsHandler"})
	groups, err := w.repository.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("query workout groups for trainerUUID: %s failed: %v", trainerUUID, err)
		return WorkoutGroupsDetails{}, ErrRepositoryFailure
	}
	out := WorkoutGroupsDetails{
		WorkoutGroups: ConvertToWorkoutGroupsDetails(groups...),
	}
	return out, nil
}

func NewWorkoutGroupsHandler(t TrainerRepository) *WorkoutGroupsHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &WorkoutGroupsHandler{
		repository: t,
	}
}
