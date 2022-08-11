package command

import (
	"context"
	"notification-service/internal/domain"
	"time"
)

type InsertNotificationHandlerRepository interface {
	InsertNotification(ctx context.Context, n domain.NotificationMessage) error
}

type InsertNotificationHandler struct {
	repo InsertNotificationHandlerRepository
}

type InsertNotificationCommand struct {
	UserUUID     string
	TrainingUUID string
	Title        string
	Content      string
	Trainer      string
	Date         string
}

func (i *InsertNotificationHandler) Do(ctx context.Context, cmd InsertNotificationCommand) error {
	t, err := time.Parse(time.RFC3339, cmd.Date)
	if err != nil {
		return err
	}

	n, err := domain.NewNotificationMessage(domain.NotificationMessageData{
		UserUUID:     cmd.UserUUID,
		TrainingUUID: cmd.TrainingUUID,
		Title:        cmd.Title,
		Content:      cmd.Content,
		Trainer:      cmd.Trainer,
		Date:         t,
	})
	if err != nil {
		return err
	}

	err = i.repo.InsertNotification(ctx, n)
	if err != nil {
		return nil
	}
	return nil
}

func NewInsertNotificationHandler(r InsertNotificationHandlerRepository) *InsertNotificationHandler {
	if r == nil {
		panic("nil insert notification repository")
	}
	h := InsertNotificationHandler{
		repo: r,
	}
	return &h
}
