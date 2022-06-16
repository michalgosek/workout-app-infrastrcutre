package command

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutRegistration struct {
	CustomerUUID string
	GroupUUID    string
}

type WorkoutDayAssigner interface {
	QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error)
	UpsertWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error
	UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error
}

type AssignCustomerWorkoutHandler struct {
	repository WorkoutDayAssigner
}

func (a *AssignCustomerWorkoutHandler) Do(ctx context.Context, args WorkoutRegistration) error {
	//logger := logrus.WithFields(logrus.Fields{"Component": "AssignCustomerHandler"})

	return nil
}

func NewAssignCustomerHandler(c WorkoutDayAssigner) *AssignCustomerWorkoutHandler {
	if c == nil {
		panic("nil repository")
	}
	return &AssignCustomerWorkoutHandler{
		repository: c,
	}
}
