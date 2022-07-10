package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	domain "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	GetCustomerWorkoutGroup(ctx context.Context, args trainer.WorkoutGroupWithCustomerArgs) (domain.WorkoutGroup, error)
}
