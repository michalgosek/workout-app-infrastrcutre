package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type AssignParticipantRepository interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
	UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error
}

type AssignParticipantHandler struct {
	repo AssignParticipantRepository
}

type AssignParticipant struct {
	TrainerUUID  string
	TrainingUUID string
	Participant  trainings.Participant
}

func (a *AssignParticipantHandler) Do(ctx context.Context, cmd AssignParticipant) error {
	training, err := a.repo.QueryTrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}
	err = training.AssignParticipant(cmd.Participant)
	if err != nil {
		return err
	}

	err = a.repo.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func NewAssignParticipantHandler(r AssignParticipantRepository) *AssignParticipantHandler {
	if r == nil {
		panic("nil assign participant repository")
	}
	h := AssignParticipantHandler{repo: r}
	return &h
}
