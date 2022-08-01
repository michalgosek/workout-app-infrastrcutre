package query

import (
	"context"
)

type TrainingsService interface {
	TrainingGroup(ctx context.Context, q TrainingQuery) (TrainingGroup, error)
}
