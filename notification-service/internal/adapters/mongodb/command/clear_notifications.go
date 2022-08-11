package command

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ClearNotificationsHandlerConfig struct {
	Database       string
	CommandTimeout time.Duration
}

type ClearNotificationsHandler struct {
	cfg ClearNotificationsHandlerConfig
	cli *mongo.Client
}

func (c *ClearNotificationsHandler) ClearNotifications(ctx context.Context, userUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	db := c.cli.Database(c.cfg.Database)
	name := NewUserCollectionName(userUUID)
	coll := db.Collection(name)
	err := coll.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewClearNotificationsHandler(cli *mongo.Client, cfg ClearNotificationsHandlerConfig) *ClearNotificationsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := ClearNotificationsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
