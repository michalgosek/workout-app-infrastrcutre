package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
)

type AssignParticipantService interface {
	AssignParticipant(ctx context.Context, groupUUID, trainerUUID string, p trainings.Participant) error
}

type AssignParticipantHandler struct {
	service AssignParticipantService
}

type AssignParticipant struct {
	TrainerUUID string
	GroupUUID   string
	Participant trainings.Participant
}

func (a *AssignParticipantHandler) Do(ctx context.Context, cmd AssignParticipant) error {
	err := a.service.AssignParticipant(ctx, cmd.GroupUUID, cmd.TrainerUUID, cmd.Participant)
	if err != nil {
		return err
	}
	return nil
}

func NewAssignParticipantHandler(s AssignParticipantService) *AssignParticipantHandler {
	if s == nil {
		panic("nil assign participant service")
	}
	h := AssignParticipantHandler{service: s}
	return &h
}
