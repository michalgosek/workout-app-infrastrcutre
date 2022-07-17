package query

import (
	"context"
)

type TrainingGroupReadModel interface {
	TrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) (TrainerWorkoutGroup, error)
}

type TrainingGroupHandler struct {
	read TrainingGroupReadModel
}

func (t *TrainingGroupHandler) Do(ctx context.Context, trainingUUID, trainerUUID string) (TrainerWorkoutGroup, error) {
	g, err := t.read.TrainingGroup(ctx, trainingUUID, trainerUUID)
	if err != nil {
		return TrainerWorkoutGroup{}, nil
	}
	return g, nil
}

func NewTrainingGroupHandler(r TrainingGroupReadModel) *TrainingGroupHandler {
	if r == nil {
		panic("nil training group read model")
	}
	h := TrainingGroupHandler{read: r}
	return &h
}
