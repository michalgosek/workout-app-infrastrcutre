package application

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
)

type Commands struct {
	ScheduleTrainerWorkoutGroup         *command.ScheduleTrainerWorkoutGroupHandler
	AssignParticipantToWorkoutGroup     *command.AssignParticipantHandler
	UnassignParticipantFromWorkoutGroup *command.UnassignParticipantHandler
	CancelTrainerWorkoutGroup           *command.CancelWorkoutGroupHandler
}

type Queries struct {
	TrainerWorkoutGroup  *query.TrainerWorkoutGroupHandler
	TrainerWorkoutGroups *query.TrainerWorkoutGroupsHandler
}

type Application struct {
	Commands
	Queries
}
