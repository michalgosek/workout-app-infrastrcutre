package command

import (
	"context"
)

type ClearNotificationsHandlerRepository interface {
	ClearNotifications(ctx context.Context, userUUID string) error
}

type ClearNotificationsHandler struct {
	repo ClearNotificationsHandlerRepository
}

func (c *ClearNotificationsHandler) Do(ctx context.Context, userUUID string) error {
	err := c.repo.ClearNotifications(ctx, userUUID)
	if err != nil {
		return nil
	}
	return nil
}

func NewClearNotificationsHandler(r ClearNotificationsHandlerRepository) *ClearNotificationsHandler {
	if r == nil {
		panic("nil clear notifications repository")
	}
	h := ClearNotificationsHandler{
		repo: r,
	}
	return &h
}
