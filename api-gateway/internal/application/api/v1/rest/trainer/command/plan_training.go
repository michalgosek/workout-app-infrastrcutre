package command

import (
	"context"
)

type PlanTrainingHandler struct {
	trainings TrainingsService
}

func (p *PlanTrainingHandler) Do(ctx context.Context, cmd PlanTrainingCommand) error {
	err := p.trainings.PlanTraining(ctx, PlanTrainingCommand{
		User: User{
			UUID: cmd.User.UUID,
			Name: cmd.User.Name,
			Role: cmd.User.Role,
		},
		GroupName: cmd.GroupName,
		GroupDesc: cmd.GroupDesc,
		Date:      cmd.Date,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewPlanTrainingHandler(t TrainingsService) *PlanTrainingHandler {
	if t == nil {
		panic("nil trainings service")
	}
	h := PlanTrainingHandler{
		trainings: t,
	}
	return &h
}
