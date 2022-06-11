package application

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

type TrainerWorkoutSessionRepoistory interface {
	UpsertTrainerWorkoutSession(ctx context.Context, s domain.TrainerWorkoutSession)
	QueryTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error)
	QueryTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error)
	CancelTrainerWorkoutSession(ctx context.Context, sessionUUID string) (domain.TrainerWorkoutSession, error)
	CancelTrainerWorkoutSessions(ctx context.Context, sessionUUIDs ...string) ([]domain.TrainerWorkoutSession, error)
}

type TrainerService struct {
	repository TrainerWorkoutSessionRepoistory
}

func (t *TrainerService) CreateTrainerWorkoutSession(ctx context.Context, trainerUUID string) error {
	return nil
}

func (t *TrainerService) GetTrainerWorkoutSession(ctx context.Context, trainerUUID string) (domain.TrainerWorkoutSession, error) {
	return domain.TrainerWorkoutSession{}, nil
}

func (t *TrainerService) GetTrainerWorkoutSessions(ctx context.Context, trainerUUID string) ([]domain.TrainerWorkoutSession, error) {
	return nil, nil
}

func (t *TrainerService) CancelTrainerWorkoutSession(ctx context.Context, sessionUUID, trainerUUID string) error {

	return nil
}
