package application

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
)

type Commands struct {
	CreateTrainerWorkout  *command.ScheduleWorkoutHandler
	UnassignCustomer      *command.UnassignCustomerHandler
	DeleteTrainerWorkout  *command.CancelWorkoutHandler
	DeleteTrainerWorkouts *command.CancelWorkoutsHandler
}

type Queries struct {
	GetTrainerWorkout  *query.WorkoutGroupHandler
	GetTrainerWorkouts *query.WorkoutGroupsHandler
}

type Application struct {
	Commands
	Queries
}
