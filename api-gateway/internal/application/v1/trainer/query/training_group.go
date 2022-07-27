package query

import (
	"context"
	"errors"
)

type TrainingHandler struct {
	users     UsersService
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
	user, err := t.users.User(ctx, q.UserUUID)
	if err != nil {
		return TrainingGroup{}, err
	}

	training, err := t.trainings.TrainingGroup(ctx, TrainingQuery{
		User: User{
			UUID: user.UUID,
			Role: user.Role,
			Name: user.Name,
		},
		TrainingUUID: q.TrainingUUID,
	})
	if err != nil {
		return TrainingGroup{}, err
	}
	return training, nil
}

func NewTrainingHandler(u UsersService, t TrainingsService) (*TrainingHandler, error) {
	if u == nil {
		return nil, errors.New("nil users service")
	}
	if t == nil {
		return nil, errors.New("nil trainings service")
	}

	h := TrainingHandler{
		users:     u,
		trainings: t,
	}
	return &h, nil
}
