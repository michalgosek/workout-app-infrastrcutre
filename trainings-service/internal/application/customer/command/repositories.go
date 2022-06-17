package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type CustomerRepository interface {
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, groupUUID string) error
}

type TrainerRepository interface {
	UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	QueryTrainerWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
}
