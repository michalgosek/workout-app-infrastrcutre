package application

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
)

type Commands struct {
	PlanTrainingGroup    *command.PlanTrainingGroupHandler
	AssignParticipant    *command.AssignParticipantHandler
	UnassignParticipant  *command.UnassignParticipantHandler
	CancelTrainingGroup  *command.CancelTrainingGroupHandler
	CancelTrainingGroups *command.CancelTrainingGroupsHandler
}

type Queries struct {
	TrainingGroup  *query.TrainingGroupHandler
	TrainingGroups *query.TrainingGroupsHandler
}

type Application struct {
	Commands
	Queries
}
