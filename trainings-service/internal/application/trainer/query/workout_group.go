package query

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=TrainerWorkoutGroupReadModel --case underscore --with-expecter
type TrainerWorkoutGroupReadModel interface {
	TrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (TrainerWorkoutGroup, error)
}

type TrainerWorkoutGroupHandler struct {
	read TrainerWorkoutGroupReadModel
}

func (t *TrainerWorkoutGroupHandler) Do(ctx context.Context, trainerUUID, groupUUID string) (TrainerWorkoutGroup, error) {
	logger := logrus.WithFields(logrus.Fields{"trainer-query": "WorkoutGroupHandler"})
	group, err := t.read.TrainerWorkoutGroup(ctx, trainerUUID, groupUUID)
	if err != nil {
		logger.Errorf("query - get trainer workout group failure: %s", err)
		return TrainerWorkoutGroup{}, fmt.Errorf("trainer service failure:%w", err)
	}
	return group, nil
}

func NewTrainerWorkoutGroupHandler(t TrainerWorkoutGroupReadModel) *TrainerWorkoutGroupHandler {
	if t == nil {
		panic("nil trainer workout group read model impl")
	}
	h := TrainerWorkoutGroupHandler{
		read: t,
	}
	return &h
}
