package query

import "context"

const UIFormat = "02/01/2006 15:04"

type Notification struct {
	UUID         string `json:"uuid"`
	UserUUID     string `json:"user_uuid"`
	TrainingUUID string `json:"training_uuid"`
	Title        string `json:"title"`
	Trainer      string `json:"trainer"`
	Content      string `json:"content"`
	Date         string `json:"date"`
}

type AllNotificationsHandlerRepository interface {
	AllNotifications(ctx context.Context, UUID string) ([]Notification, error)
}

type AllNotificationsHandler struct {
	repo AllNotificationsHandlerRepository
}

func (a *AllNotificationsHandler) Do(ctx context.Context, UUID string) ([]Notification, error) {
	notifications, err := a.repo.AllNotifications(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func NewAllNotificationsHandler(r AllNotificationsHandlerRepository) *AllNotificationsHandler {
	if r == nil {
		panic("nil all notifications repository")
	}
	h := AllNotificationsHandler{
		repo: r,
	}
	return &h
}
