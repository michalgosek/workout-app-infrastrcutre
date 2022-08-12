package command

import (
	"context"
)

type UnassignParticipantHandler struct {
	command UpdateTrainingGroupRepository
	query   TrainingGroupRepository
}

type UnassignParticipant struct {
	TrainingUUID    string
	TrainerUUID     string
	ParticipantUUID string
}

func (u *UnassignParticipantHandler) Do(ctx context.Context, cmd UnassignParticipant) error {
	training, err := u.query.TrainingGroup(ctx, cmd.TrainingUUID)
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
	err = u.command.UpdateTrainingGroup(ctx, &training)
	if err != nil {
		return err
	}
	return nil
}

func NewUnassignParticipantHandler(cmd UpdateTrainingGroupRepository, query TrainingGroupRepository) *UnassignParticipantHandler {
	if cmd == nil {
		panic("nil update training group repository")
	}
	if query == nil {
		panic("nil query training group repository")
	}

	h := UnassignParticipantHandler{
		command: cmd,
		query:   query,
	}
	return &h
}
