package command

import (
	"context"
)

type UnassignParticipantService interface {
	UnassignParticipant(ctx context.Context, trainingUUID, trainerUUID, participantUUID string) error
}

type UnassignParticipantHandler struct {
	service UnassignParticipantService
}

type UnassignParticipant struct {
	TrainingUUID    string
	TrainerUUID     string
	ParticipantUUID string
}

func (u *UnassignParticipantHandler) Do(ctx context.Context, cmd UnassignParticipant) error {
	err := u.service.UnassignParticipant(ctx, cmd.TrainingUUID, cmd.TrainerUUID, cmd.ParticipantUUID)
	if err != nil {
		return err
	}
	return nil
}

func NewUnassignParticipantHandler(s UnassignParticipantService) *UnassignParticipantHandler {
	if s == nil {
		panic("nil unassign participant service")
	}
	h := UnassignParticipantHandler{service: s}
	return &h
}
