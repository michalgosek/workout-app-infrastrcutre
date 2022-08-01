package command

import (
	"context"
)

type TrainingsService interface {
	PlanTraining(ctx context.Context, cmd PlanTrainingCommand) error
}
