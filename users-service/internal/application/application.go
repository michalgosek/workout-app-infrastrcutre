package application

import (
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/command"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
)

type Commands struct {
	RegisterUser *command.RegisterUserHandler
}

type Queries struct {
	User *query.UserHandler
}

type Application struct {
	Commands
	Queries
}
