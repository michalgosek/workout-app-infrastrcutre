package query

import (
	"context"
)

type TrainingHandler struct {
	trainings TrainingsService
}

type Training struct {
	TrainingUUID string
	UserUUID     string
}

type TrainingQuery struct {
	User         User
	TrainingUUID string
}

func (t *TrainingHandler) Do(ctx context.Context, q Training) (TrainingGroup, error) {
	training, err := t.trainings.TrainingGroup(ctx, TrainingQuery{
		TrainingUUID: q.TrainingUUID,
	})
	if err != nil {
		return TrainingGroup{}, err
	}
	return training, nil
}

func NewTrainingHandler(t TrainingsService) *TrainingHandler {
	if t == nil {
		panic("nil trainings service")
	}

	h := TrainingHandler{
		trainings: t,
	}
	return &h
}
