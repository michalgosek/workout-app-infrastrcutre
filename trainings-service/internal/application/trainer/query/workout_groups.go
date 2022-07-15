package query

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=TrainerWorkoutGroupsReadModel --case underscore --with-expecter
type TrainerWorkoutGroupsReadModel interface {
	TrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error)
}

type TrainerWorkoutGroupsHandler struct {
	model TrainerWorkoutGroupsReadModel
}

func (t *TrainerWorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"trainer-query": "GetWorkoutsHandler"})
	groups, err := t.model.TrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("query - get trainer workout group failure: %s", err)
		return nil, fmt.Errorf("trainer service failure:%w", err)
	}
	return groups, nil
}

func NewTrainerWorkoutGroupsHandler(t TrainerWorkoutGroupsReadModel) *TrainerWorkoutGroupsHandler {
	if t == nil {
		panic("nil trainer workout group read model")
	}
	h := TrainerWorkoutGroupsHandler{
		model: t,
	}
	return &h
}
