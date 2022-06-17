package query

import (
	"context"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
)

type CustomerWorkoutDay struct {
	Date         time.Time
	Trainer      string
	WorkoutName  string
	WorkoutDesc  string
	Participants int
}

type WorkoutDayHandlerRepository interface {
	QueryCustomerWorkoutDay(ctx context.Context, customerUUID, groupUUID string) (customer.WorkoutDay, error)
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
}

type WorkoutDayHandler struct {
	repository WorkoutDayHandlerRepository
}

type WorkoutDayDetails struct {
	CustomerUUID string
	GroupUUID    string
}

func (w *WorkoutDayHandler) Do(ctx context.Context, d WorkoutDayDetails) (CustomerWorkoutDay, error) {

	return CustomerWorkoutDay{}, nil
}

func NewWorkoutDayHandler(r WorkoutDayHandlerRepository) *WorkoutDayHandler {
	return &WorkoutDayHandler{repository: r}
}
