package command

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type PlanTrainingGroup struct {
	UUID        string
	Name        string
	Description string
	Date        time.Time
	Trainer     trainings.Trainer
}

type CreateTrainingGroupRepository interface {
	InsertTrainerGroup(ctx context.Context, g *trainings.TrainingGroup) error
	IsTrainingGroupDuplicated(ctx context.Context, g *trainings.TrainingGroup) (bool, error)
}

type PlanTrainingGroupHandler struct {
	repo CreateTrainingGroupRepository
}

func (p *PlanTrainingGroupHandler) Do(ctx context.Context, cmd PlanTrainingGroup) (string, error) {
	g, err := trainings.NewTrainingGroup(cmd.UUID, cmd.Name, cmd.Description, cmd.Date, cmd.Trainer)
	if err != nil {
		return "", err
	}
	duplicate, err := p.repo.IsTrainingGroupDuplicated(ctx, g)
	if duplicate {
		return "", ErrTrainingDuplicated
	}
	err = p.repo.InsertTrainerGroup(ctx, g)
	if err != nil {
		return "", err
	}

	return g.UUID(), nil
}

func NewPlanTrainingGroupHandler(r CreateTrainingGroupRepository) *PlanTrainingGroupHandler {
	if r == nil {
		panic("nil create training group repository")
	}
	h := PlanTrainingGroupHandler{repo: r}
	return &h
}

var ErrTrainingDuplicated = errors.New("training group duplicated")
