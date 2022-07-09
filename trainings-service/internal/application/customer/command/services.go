package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services"
)

//go:generate mockery --name=CustomerService --case underscore --with-expecter
type CustomerService interface {
	CancelWorkoutDay(ctx context.Context, args services.CancelWorkoutDayArgs) error
	ScheduleWorkoutDay(ctx context.Context, args services.ScheduleWorkoutDayArgs) error
}

//go:generate mockery --name=TrainerService --case underscore --with-expecter
type TrainerService interface {
	CancelCustomerWorkoutParticipation(ctx context.Context, args services.CancelCustomerWorkoutParticipationArgs) error
	AssignCustomerToWorkoutGroup(ctx context.Context, args services.AssignCustomerToWorkoutGroupArgs) (services.AssignedCustomerWorkoutGroupDetails, error)
}
