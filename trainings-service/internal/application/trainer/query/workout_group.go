package query

import (
	"context"

	"github.com/sirupsen/logrus"
)

type WorkoutGroupHandler struct {
	repository TrainerRepository
}

func (t *WorkoutGroupHandler) Do(ctx context.Context, groupUUID, trainerUUID string) (WorkoutGroupDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "WorkoutGroupHandler"})

	group, err := t.repository.QueryTrainerWorkoutGroup(ctx, groupUUID)
	if err != nil {
		logger.Errorf("query workout groupUUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return WorkoutGroupDetails{}, ErrRepositoryFailure
	}
	if group.UUID() == "" {
		return WorkoutGroupDetails{}, nil
	}
	if group.TrainerUUID() != trainerUUID {
		logger.Errorf("query workout group UUID: %s does not belong to trainerUUID: %s", groupUUID, trainerUUID)
		return WorkoutGroupDetails{}, ErrWorkoutGroupNotOwner
	}

	out := WorkoutGroupDetails{
		TrainerUUID: group.TrainerUUID(),
		TrainerName: group.TrainerName(),
		GroupUUID:   group.UUID(),
		GroupDesc:   group.Description(),
		GroupName:   group.Name(),
		Date:        group.Date().String(),
	}
	return out, nil
}

func NewWorkoutGroupHandler(t TrainerRepository) *WorkoutGroupHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &WorkoutGroupHandler{
		repository: t,
	}
}
