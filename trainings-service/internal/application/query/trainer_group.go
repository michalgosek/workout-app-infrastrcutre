package query

import (
	"context"
)

type TrainingGroupReadModel interface {
	TrainerGroup(ctx context.Context, trainingUUID, trainerUUID string) (TrainerWorkoutGroup, error)
}

type TrainerGroupHandler struct {
	read TrainingGroupReadModel
}

func (t *TrainerGroupHandler) Do(ctx context.Context, trainingUUID, trainerUUID string) (TrainerWorkoutGroup, error) {
	g, err := t.read.TrainerGroup(ctx, trainingUUID, trainerUUID)
	if err != nil {
		return TrainerWorkoutGroup{}, nil
	}
	return g, nil
}

func NewTrainerGroupHandler(r TrainingGroupReadModel) *TrainerGroupHandler {
	if r == nil {
		panic("nil trainer group read model")
	}
	h := TrainerGroupHandler{read: r}
	return &h
}
