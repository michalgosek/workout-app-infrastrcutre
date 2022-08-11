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

type CancelTrainingGroupRepository interface {
	QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error)
	DeleteTrainerGroup(ctx context.Context, trainingUUID, trainerUUID string) error
}

type NotificationService interface {
	CreateNotification(n Notification) error
}

type CancelTrainingGroupHandler struct {
	repo    CancelTrainingGroupRepository
	service NotificationService
}

type CancelWorkoutGroup struct {
	TrainingUUID string
	TrainerUUID  string
}

func (c *CancelTrainingGroupHandler) Do(ctx context.Context, cmd CancelWorkoutGroup) error {
	training, err := c.repo.QueryTrainingGroup(ctx, cmd.TrainingUUID)
	if err != nil {
		return err
	}
	if !training.IsOwnedByTrainer(cmd.TrainerUUID) {
		return ErrTrainingNotOwnedByTrainer
	}

	err = c.repo.DeleteTrainerGroup(ctx, cmd.TrainingUUID, cmd.TrainerUUID)
	if err != nil {
		return err
	}
	for _, p := range training.Participants() {
		err := c.service.CreateNotification(Notification{
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

func NewCancelTrainingGroupHandler(r CancelTrainingGroupRepository, s NotificationService) *CancelTrainingGroupHandler {
	if r == nil {
		panic("nil cancel training group repository")
	}
	if s == nil {
		panic("nil notification service")
	}
	h := CancelTrainingGroupHandler{
		repo:    r,
		service: s,
	}
	return &h
}
