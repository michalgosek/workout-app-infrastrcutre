package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type AssignParticipantHandler struct {
	command UpdateTrainingGroupRepository
	query   TrainingGroupRepository
}

type AssignParticipant struct {
	TrainerUUID  string
	TrainingUUID string
	Participant  trainings.Participant
}

func (a *AssignParticipantHandler) Do(ctx context.Context, cmd AssignParticipant) error {
	training, err := a.query.TrainingGroup(ctx, cmd.TrainingUUID)
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

	err = a.command.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func NewAssignParticipantHandler(command UpdateTrainingGroupRepository, query TrainingGroupRepository) *AssignParticipantHandler {
	if command == nil {
		panic("nil update training group repository")
	}
	if query == nil {
		panic("nil query training group repository")
	}
	h := AssignParticipantHandler{
		command: command,
		query:   query,
	}
	return &h
}
