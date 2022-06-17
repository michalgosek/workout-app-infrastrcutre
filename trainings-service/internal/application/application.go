package application

import (
	custcmd "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/command"
	trainercmd "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	trainerqry "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
)

type Commands struct {
	CreateTrainerWorkout    *trainercmd.ScheduleWorkoutHandler
	UnassignCustomer        *trainercmd.UnassignCustomerHandler
	DeleteTrainerWorkout    *trainercmd.CancelWorkoutHandler
	DeleteTrainerWorkouts   *trainercmd.CancelWorkoutsHandler
	CustomerScheduleWorkout *custcmd.ScheduleWorkoutHandler
}

type Queries struct {
	GetTrainerWorkout  *trainerqry.WorkoutGroupHandler
	GetTrainerWorkouts *trainerqry.WorkoutGroupsHandler
}

type Application struct {
	Commands
	Queries
}
