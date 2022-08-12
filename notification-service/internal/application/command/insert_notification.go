package command

import (
	"context"
	"github.com/sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{"InsertNotificationHandler": "Do"})

	t, err := time.Parse(time.RFC3339, cmd.Date)
	if err != nil {
		logrus.WithFields(logrus.Fields{"TrainingUUID": cmd.TrainingUUID, "UserUUID": cmd.UserUUID}).Errorf("notification message time parse failure %s", err)
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
		logrus.WithFields(logrus.Fields{"UserUUID": n.UserUUID(), "Date": n.Date()}).Errorf("create notification message failure: %s", err)
		return err
	}

	err = i.repo.InsertNotification(ctx, n)
	if err != nil {
		logrus.WithFields(logrus.Fields{"UserUUID": n.UserUUID(), "Date": n.Date()}).Errorf("insert notification message failure: %s", err)
		return nil
	}

	logrus.WithFields(logrus.Fields{"MessageUUID": n.UUID(), "UserUUID": n.UserUUID(), "Date": n.Date()}).Info("insert notification message success")
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
