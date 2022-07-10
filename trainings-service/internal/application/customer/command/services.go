package command

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
)

//go:generate mockery --name=TrainingsService --case underscore --with-expecter
type TrainingsService interface {
	CancelCustomerWorkout(ctx context.Context, args trainings.CancelCustomerWorkoutArgs) error
	AssignCustomerToWorkoutGroup(ctx context.Context, args trainings.AssignCustomerToWorkoutArgs) error
}

var ErrNilTrainingsService = errors.New("nil trainings service")
