package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"time"
)

type CustomerRepository interface {
	QueryCustomerWorkoutDay(ctx context.Context, customerUUID, workoutDayUUID string) (customer.WorkoutDay, error)
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, workoutDayUUID string) error
}

type TrainerRepository interface {
	UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	QueryTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) (trainer.WorkoutGroup, error)
	QueryTrainerWorkoutGroupWithDate(ctx context.Context, trainerUUID string, date time.Time) (trainer.WorkoutGroup, error)
	DeleteTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) error
	DeleteTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error
}
