package application

import (
	command2 "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
)

type Commands struct {
	CreateTrainerWorkout  *command.CreateWorkoutHandler
	AssignCustomer        *command2.AssignCustomerHandler
	UnassignCustomer      *command.UnassignCustomerHandler
	DeleteTrainerWorkout  *command.WorkoutDeleteHandler
	DeleteTrainerWorkouts *command.WorkoutsDeleteHandler
}

type Queries struct {
	GetTrainerWorkout  *query.GetWorkoutHandler
	GetTrainerWorkouts *query.GetWorkoutsHandler
}

type Application struct {
	Commands
	Queries
}
