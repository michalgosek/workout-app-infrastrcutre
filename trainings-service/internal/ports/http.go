package ports

import (
	"context"
	"net/http"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type TrainingService interface {
	CreateTrainerSchedule(ctx context.Context, args application.TrainerSchedule) error
	GetSchedule(ctx context.Context, scheduleUUID, trainerUUID string) (trainer.TrainerSchedule, error)
	GetSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error)
	AssingCustomer(ctx context.Context, customerUUID, scheduleUUID, trainerUUID string) error
	DeleteSchedule(ctx context.Context, sessionUUID, trainerUUID string) error
	DeleteSchedules(ctx context.Context, sessionUUID string) error
}
type HTTP struct {
	service TrainingService
}

func (h *HTTP) CreateTrainerSchedule(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTP) GetSchedule(w http.ResponseWriter, r *http.Request) {

}

func NewHTTP(service TrainingService) *HTTP {
	h := HTTP{
		service: service,
	}
	return &h
}
