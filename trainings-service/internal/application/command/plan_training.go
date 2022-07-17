package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type PlanTrainingGroup struct {
	UUID        string
	Name        string
	Description string
	Date        time.Time
	TrainerUUID string
	TrainerName string
}

type CreateTrainingGroupService interface {
	CreateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
}

type PlanTrainingGroupHandler struct {
	service CreateTrainingGroupService
}

func (p *PlanTrainingGroupHandler) Do(ctx context.Context, cmd PlanTrainingGroup) (string, error) {
	t, err := trainings.NewTrainer(cmd.TrainerUUID, cmd.TrainerName)
	if err != nil {
		return "", err
	}
	g, err := trainings.NewTrainingGroup(cmd.UUID, cmd.Name, cmd.Description, cmd.Date, t)
	if err != nil {
		return "", err
	}

	err = p.service.CreateTrainingGroup(ctx, g)
	if err != nil {
		return "", err
	}
	return g.UUID(), nil
}

func NewPlanTrainingGroupHandler(s CreateTrainingGroupService) *PlanTrainingGroupHandler {
	if s == nil {
		panic("nil create training group service")
	}
	h := PlanTrainingGroupHandler{service: s}
	return &h
}
