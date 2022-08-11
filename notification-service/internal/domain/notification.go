package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"notification-service/internal/adapters/mongodb/documents"
	"time"
)

type NotificationMessage struct {
	uuid         string
	userUUID     string
	trainingUUID string
	trainer      string
	title        string
	content      string
	date         time.Time
}

func (n NotificationMessage) UUID() string {
	return n.uuid
}

func (n NotificationMessage) UserUUID() string {
	return n.userUUID
}

func (n NotificationMessage) TrainingUUID() string {
	return n.trainingUUID
}

func (n NotificationMessage) Trainer() string {
	return n.trainer
}

func (n NotificationMessage) Date() time.Time {
	return n.date
}

func (n NotificationMessage) Title() string {
	return n.title
}

func (n NotificationMessage) Content() string {
	return n.content
}

type NotificationMessageData struct {
	UserUUID     string
	TrainingUUID string
	Title        string
	Content      string
	Trainer      string
	Date         time.Time
}

func (n *NotificationMessageData) Validate() error {
	if n.Date.IsZero() {
		return errors.New("zero value for notification")
	}
	m := map[string]string{
		"user uuid":     n.UserUUID,
		"training uuid": n.TrainingUUID,
		"title":         n.Title,
		"trainer":       n.Trainer,
		"content":       n.Content,
	}
	for k, v := range m {
		if v == "" {
			text := fmt.Sprintf("missing value for param %s", k)
			return errors.New(text)
		}
	}
	return nil
}

func NewNotificationMessage(d NotificationMessageData) (NotificationMessage, error) {
	err := d.Validate()
	if err != nil {
		return NotificationMessage{}, nil
	}
	n := NotificationMessage{
		uuid:         uuid.NewString(),
		userUUID:     d.UserUUID,
		trainingUUID: d.TrainingUUID,
		trainer:      d.Trainer,
		title:        d.Title,
		content:      d.Content,
		date:         d.Date,
	}
	return n, nil
}

func ConvertToDomainNotification(doc documents.NotificationWriteModel) NotificationMessage {
	return NotificationMessage{
		uuid:         doc.UUID,
		userUUID:     doc.UserUUID,
		trainingUUID: doc.TrainingUUID,
		title:        doc.Title,
		content:      doc.Content,
		trainer:      doc.Trainer,
		date:         doc.Date,
	}
}
