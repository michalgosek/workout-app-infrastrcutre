package command

import (
	"context"
	"fmt"
	"time"
)

type TrainerParticipant struct {
	TrainingUUID string
	TrainingName string
	Trainer      string
	UserUUID     string
	Date         time.Time
}

type DeleteTrainingGroupsRepository interface {
	DeleteTrainingGroups(ctx context.Context, trainerUUID string) error
}

type TrainerParticipantsRepository interface {
	TrainerParticipants(ctx context.Context, trainerUUID string) ([]TrainerParticipant, error)
}

type CancelTrainingGroupsHandler struct {
	command DeleteTrainingGroupsRepository
	query   TrainerParticipantsRepository
	service NotificationService
}

func (c *CancelTrainingGroupsHandler) Do(ctx context.Context, trainerUUID string) error {
	participants, err := c.query.TrainerParticipants(ctx, trainerUUID)
	if err != nil {
		return err
	}

	err = c.command.DeleteTrainingGroups(ctx, trainerUUID)
	if err != nil {
		return err
	}
	for _, p := range participants {
		err = c.service.CreateNotification(ctx, Notification{
			UserUUID:     p.UserUUID,
			TrainingUUID: p.TrainingUUID,
			Title:        fmt.Sprintf("Training canceled - %s", p.TrainingName),
			Content:      fmt.Sprintf("Training session has been canceled"),
			Trainer:      p.Trainer,
			Date:         p.Date,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewCancelTrainingGroupsHandler(command DeleteTrainingGroupsRepository, query TrainerParticipantsRepository, n NotificationService) *CancelTrainingGroupsHandler {
	if command == nil {
		panic("nil delete training groups repository")
	}
	if query == nil {
		panic("nil delete trainer participants repository")
	}
	if n == nil {
		panic("nil notification service")
	}
	h := CancelTrainingGroupsHandler{
		command: command,
		query:   query,
		service: n,
	}

	return &h
}
