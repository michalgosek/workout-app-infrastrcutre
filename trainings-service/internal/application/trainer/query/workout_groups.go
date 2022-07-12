package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	GetTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
}

type WorkoutGroupsDetails struct {
	WorkoutGroups []WorkoutGroupDetails
}

type WorkoutGroupsHandler struct {
	trainerService TrainerService
}

func (w *WorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) (WorkoutGroupsDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "GetWorkoutsHandler"})
	groups, err := w.trainerService.GetTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("query - get trainer workout group failure: %s", err)
		return WorkoutGroupsDetails{}, fmt.Errorf("trainer service failure:%w", err)
	}
	out := WorkoutGroupsDetails{
		WorkoutGroups: ConvertToWorkoutGroupsDetails(groups...),
	}
	return out, nil
}

func NewWorkoutGroupsHandler(t TrainerService) (*WorkoutGroupsHandler, error) {
	if t == nil {
		return nil, ErrNilTrainerService
	}
	h := WorkoutGroupsHandler{
		trainerService: t,
	}
	return &h, nil
}

var ErrNilTrainerService = errors.New("nil trainer service")
