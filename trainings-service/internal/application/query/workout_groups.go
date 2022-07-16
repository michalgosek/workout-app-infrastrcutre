package query

import "context"

type TrainerWorkoutGroupsReadModel interface {
	TrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error)
}

type TrainerWorkoutGroupsHandler struct {
	read TrainerWorkoutGroupsReadModel
}

type TrainerWorkoutGroups struct {
	TrainerUUID string
	GroupUUID   string
}

func (t *TrainerWorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) ([]TrainerWorkoutGroup, error) {
	gg, err := t.read.TrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func NewTrainerWorkoutGroupsHandler(r TrainerWorkoutGroupsReadModel) *TrainerWorkoutGroupsHandler {
	if r == nil {
		panic("nil trainer workout groups read model")
	}
	h := TrainerWorkoutGroupsHandler{read: r}
	return &h
}
