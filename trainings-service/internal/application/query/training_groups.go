package query

import "context"

type TrainerGroupsReadModel interface {
	TrainingGroups(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error)
}

type TrainingGroupsHandler struct {
	read TrainerGroupsReadModel
}

func (t *TrainingGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error) {
	gg, err := t.read.TrainingGroups(ctx, trainerUUID)
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func NewTrainingGroupsHandlerHandler(r TrainerGroupsReadModel) *TrainingGroupsHandler {
	if r == nil {
		panic("nil training groups read model")
	}
	h := TrainingGroupsHandler{read: r}
	return &h
}
