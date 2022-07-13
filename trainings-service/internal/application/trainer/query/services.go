package query

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	GetTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error)
	GetTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (trainer.WorkoutGroup, error)
}

var ErrNilTrainerService = errors.New("nil trainer service")
