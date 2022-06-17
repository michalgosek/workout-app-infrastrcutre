package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupHandler struct {
	repository TrainerRepository
}

func (t *WorkoutGroupHandler) Do(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "WorkoutGroupHandler"})

	group, err := t.repository.QueryTrainerWorkoutGroup(ctx, groupUUID)
	if err != nil {
		logger.Errorf("query workout groupUUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return trainer.WorkoutGroup{}, ErrRepositoryFailure
	}
	if group.UUID() == "" {
		return trainer.WorkoutGroup{}, nil
	}
	if group.TrainerUUID() != trainerUUID {
		logger.Errorf("query workout group UUID: %s does not belong to trainerUUID: %s", groupUUID, trainerUUID)
		return trainer.WorkoutGroup{}, ErrWorkoutGroupNotOwner
	}
	return group, nil
}

func NewWorkoutGroupHandler(t TrainerRepository) *WorkoutGroupHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &WorkoutGroupHandler{
		repository: t,
	}
}
