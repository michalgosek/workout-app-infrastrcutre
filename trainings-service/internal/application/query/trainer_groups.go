package query

import "context"

type TrainerGroupsReadModel interface {
	TrainerGroups(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error)
}

type TrainerGroupsHandler struct {
	read TrainerGroupsReadModel
}

func (t *TrainerGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error) {
	gg, err := t.read.TrainerGroups(ctx, trainerUUID)
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func NewTrainerGroupsHandler(r TrainerGroupsReadModel) *TrainerGroupsHandler {
	if r == nil {
		panic("nil trainer groups read model")
	}
	h := TrainerGroupsHandler{read: r}
	return &h
}
