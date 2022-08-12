package command

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"log"
	"time"
)

type Notification struct {
	UserUUID     string    `json:"user_uuid"`
	TrainingUUID string    `json:"training_uuid"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Trainer      string    `json:"trainer"`
	Date         time.Time `json:"date"`
}

type TrainingGroupRepository interface {
	TrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
}

type DeleteTrainingGroupRepository interface {
	DeleteTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) error
}

type NotificationService interface {
	CreateNotification(ctx context.Context, n Notification) error
}

type CancelTrainingGroupHandler struct {
	query   TrainingGroupRepository
	command DeleteTrainingGroupRepository
	service NotificationService
}

type CancelWorkoutGroup struct {
	TrainingUUID string
	TrainerUUID  string
}

func (c *CancelTrainingGroupHandler) Do(ctx context.Context, cmd CancelWorkoutGroup) error {
	training, err := c.query.TrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = c.command.DeleteTrainingGroup(ctx, cmd.TrainingUUID, cmd.TrainerUUID)
	if err != nil {
		return err
	}
	for _, p := range training.Participants() { // send single batch as array instead of loop!!
		err := c.service.CreateNotification(ctx, Notification{
			UserUUID:     p.UUID(),
			TrainingUUID: training.UUID(),
			Title:        fmt.Sprintf("Training canceled - %s ", training.Name()),
			Content:      "Training session has been canceled",
			Trainer:      training.Trainer().Name(),
			Date:         training.Date(),
		})
		if err != nil {
			log.Print(err)
		}
	}
	return nil
}

func NewCancelTrainingGroupHandler(query TrainingGroupRepository, cmd DeleteTrainingGroupRepository, s NotificationService) *CancelTrainingGroupHandler {
	if query == nil {
		panic("nil query training group repository")
	}
	if cmd == nil {
		panic("nil delete training group repository")
	}
	if s == nil {
		panic("nil notification service")
	}
	h := CancelTrainingGroupHandler{
		command: cmd,
		query:   query,
		service: s,
	}
	return &h
}
