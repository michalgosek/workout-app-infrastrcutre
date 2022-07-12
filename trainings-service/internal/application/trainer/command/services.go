package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
)

//go:generate mockery --name=TrainingsService --case underscore --with-expecter
type TrainingsService interface {
	CancelTrainerWorkoutGroup(ctx context.Context, args trainings.CancelTrainerWorkoutGroupArgs) error
	CancelTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error
}
