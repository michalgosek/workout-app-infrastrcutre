package query

import (
	"context"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupHandler struct {
	repository TrainerRepository
}

func (t *WorkoutGroupHandler) Do(ctx context.Context, trainerUUID, groupUUID string) (WorkoutGroupDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "WorkoutGroupHandler"})
	group, err := t.repository.QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID)
	if err != nil {
		logger.Errorf("query workout groupUUID: %s for trainerUUID: %s failed, reason: %v", groupUUID, trainerUUID, err)
		return WorkoutGroupDetails{}, ErrRepositoryFailure
	}
	if group.UUID() == "" {
		return WorkoutGroupDetails{}, nil
	}
	out := WorkoutGroupDetails{
		TrainerUUID: group.TrainerUUID(),
		TrainerName: group.TrainerName(),
		GroupUUID:   group.UUID(),
		GroupDesc:   group.Description(),
		GroupName:   group.Name(),
		Date:        group.Date().String(),
		Customers:   convertToCustomersData(group.CustomerDetails()),
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
