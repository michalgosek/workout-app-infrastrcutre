package query

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupArgs struct {
	TrainerUUID string
	GroupUUID   string
}

type WorkoutGroupHandler struct {
	trainerService TrainerService
}

func (t *WorkoutGroupHandler) Do(ctx context.Context, args WorkoutGroupArgs) (WorkoutGroupDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"trainer-query": "WorkoutGroupHandler"})
	group, err := t.trainerService.GetTrainerWorkoutGroup(ctx, args.TrainerUUID, args.GroupUUID)
	if err != nil {
		logger.Errorf("query - get trainer workout group failure: %s", err)
		return WorkoutGroupDetails{}, fmt.Errorf("trainer service failure:%w", err)
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

func NewWorkoutGroupHandler(t TrainerService) (*WorkoutGroupHandler, error) {
	if t == nil {
		return nil, ErrNilTrainerService
	}
	h := WorkoutGroupHandler{
		trainerService: t,
	}
	return &h, nil
}
