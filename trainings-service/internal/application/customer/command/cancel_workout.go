package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type CancelWorkoutHandlerRepository interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
	DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, customerWorkoutDayUUID string) error
}

type CancelWorkoutHandler struct {
	repository CancelWorkoutHandlerRepository
}

type CancelWorkoutDetails struct {
	CustomerUUID string
	GroupUUID    string
}

func (c *CancelWorkoutHandler) Do(ctx context.Context, w CancelWorkoutDetails) error {
	return nil
}

func NewCancelWorkoutHandler(r CancelWorkoutHandlerRepository) *CancelWorkoutHandler {
	return &CancelWorkoutHandler{repository: r}
}
