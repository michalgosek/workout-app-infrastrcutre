package command

import (
	"context"
	"errors"
)

type PlanTrainingHandler struct {
	trainings TrainingsService
	users     UsersService
}

func (p *PlanTrainingHandler) Do(ctx context.Context, r PlanTraining) error {
	user, err := p.users.User(ctx, r.UserUUID)
	if err != nil {
		return err
	}
	err = p.trainings.PlanTraining(ctx, PlanTrainingCommand{
		User: User{
			UUID: user.UUID,
			Name: user.Name,
			Role: user.Role,
		},
		GroupName: r.GroupName,
		GroupDesc: r.GroupDesc,
		Date:      r.Date,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewPlanTrainingHandler(u UsersService, t TrainingsService) (*PlanTrainingHandler, error) {
	if u == nil {
		return nil, errors.New("nil user service")
	}
	if t == nil {
		return nil, errors.New("nil trainings service")
	}
	h := PlanTrainingHandler{
		trainings: t,
		users:     u,
	}
	return &h, nil
}
