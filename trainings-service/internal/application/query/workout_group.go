package query

import (
	"context"
)

type TrainerWorkoutGroupReadModel interface {
	TrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (TrainerWorkoutGroup, error)
}

type TrainerWorkoutGroupHandler struct {
	read TrainerWorkoutGroupReadModel
}

func (t *TrainerWorkoutGroupHandler) Do(ctx context.Context, groupUUID, trainerUUID string) (TrainerWorkoutGroup, error) {
	g, err := t.read.TrainerWorkoutGroup(ctx, groupUUID, trainerUUID)
	if err != nil {
		return TrainerWorkoutGroup{}, nil
	}
	return g, nil
}

func NewTrainerWorkoutGroupHandler(r TrainerWorkoutGroupReadModel) *TrainerWorkoutGroupHandler {
	if r == nil {
		panic("nil trainer workout group read model")
	}
	h := TrainerWorkoutGroupHandler{read: r}
	return &h
}
