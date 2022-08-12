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

type InsertTrainingGroupRepository interface {
	InsertTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
}

type IsTrainingGroupExistsRepository interface {
	IsTrainingGroupExists(ctx context.Context, g *trainings.TrainingGroup) (bool, error)
}

type PlanTrainingGroupHandler struct {
	command InsertTrainingGroupRepository
	query   IsTrainingGroupExistsRepository
}

func (p *PlanTrainingGroupHandler) Do(ctx context.Context, cmd PlanTrainingGroup) (string, error) {
	g, err := trainings.NewTrainingGroup(cmd.UUID, cmd.Name, cmd.Description, cmd.Date, cmd.Trainer)
	if err != nil {
		return "", err
	}
	duplicate, err := p.query.IsTrainingGroupExists(ctx, g)
	if duplicate {
		return "", ErrTrainingDuplicated
	}
	err = p.command.InsertTrainingGroup(ctx, g)
	if err != nil {
		return "", err
	}

	return g.UUID(), nil
}

func NewPlanTrainingGroupHandler(cmd InsertTrainingGroupRepository, query IsTrainingGroupExistsRepository) *PlanTrainingGroupHandler {
	if cmd == nil {
		panic("nil insert training group repository")
	}
	if query == nil {
		panic("nil query training group repository")
	}
	h := PlanTrainingGroupHandler{command: cmd, query: query}
	return &h
}

var ErrTrainingDuplicated = errors.New("training group duplicated")
