package ports

import "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"

type HTTP struct {
	service application.TrainerService
}
