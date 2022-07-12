package command

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
)

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	CreateWorkoutGroup(ctx context.Context, args trainer.CreateWorkoutGroupArgs) error
}

//go:generate mockery --name=TrainingsService --case underscore --with-expecter
type TrainingsService interface {
	CancelTrainerWorkoutGroup(ctx context.Context, args trainings.CancelTrainerWorkoutGroupArgs) error
	CancelTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error
	CancelCustomerWorkout(ctx context.Context, args trainings.CancelCustomerWorkoutArgs) error
}

var ErrNilTrainerService = errors.New("nil trainer service")
