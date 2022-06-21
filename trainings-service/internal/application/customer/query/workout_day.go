package query

import (
	"context"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
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

type WorkoutDay struct {
	CustomerUUID string
	GroupUUID    string
}

func (w *WorkoutDayHandler) Do(ctx context.Context, d WorkoutDay) (CustomerWorkoutDay, error) {
	group, err := w.repository.QueryWorkoutGroup(ctx, d.GroupUUID)
	if err != nil {
		return CustomerWorkoutDay{}, err
	}
	day, err := w.repository.QueryCustomerWorkoutDay(ctx, d.CustomerUUID, group.UUID())
	if err != nil {
		return CustomerWorkoutDay{}, err
	}
	out := CustomerWorkoutDay{
		Date:         day.Date(),
		Trainer:      group.TrainerName(),
		WorkoutName:  group.Name(),
		WorkoutDesc:  group.Description(),
		Participants: group.AssignedCustomers(),
	}
	return out, nil
}

func NewWorkoutDayHandler(r WorkoutDayHandlerRepository) *WorkoutDayHandler {
	return &WorkoutDayHandler{repository: r}
}
