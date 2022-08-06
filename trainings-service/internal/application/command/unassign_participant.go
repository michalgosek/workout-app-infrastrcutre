package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type UnassignParticipantRepository interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
	UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
}

type UnassignParticipantHandler struct {
	repo UnassignParticipantRepository
}

type UnassignParticipant struct {
	TrainingUUID    string
	TrainerUUID     string
	ParticipantUUID string
}

func (u *UnassignParticipantHandler) Do(ctx context.Context, cmd UnassignParticipant) error {
	training, err := u.repo.QueryTrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = training.UnassignParticipant(cmd.ParticipantUUID)
	if err != nil {
		return err
	}
	err = u.repo.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func NewUnassignParticipantHandler(r UnassignParticipantRepository) *UnassignParticipantHandler {
	if r == nil {
		panic("nil unassign participant repository")
	}
	h := UnassignParticipantHandler{repo: r}
	return &h
}
