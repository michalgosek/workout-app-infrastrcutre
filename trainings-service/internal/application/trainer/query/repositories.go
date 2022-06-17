package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainerRepository interface {
	QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	QueryTrainerWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
}
