package command

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-service/internal/adapters/mongodb/documents"
	"notification-service/internal/domain"
	"time"
)

func NewUserCollectionName(UUID string) string {
	return fmt.Sprintf("%s_user", UUID)
}

type InsertNotificationHandlerConfig struct {
	Database       string
	CommandTimeout time.Duration
}

type InsertNotificationHandler struct {
	cfg InsertNotificationHandlerConfig
	cli *mongo.Client
}

func (i *InsertNotificationHandler) InsertNotification(ctx context.Context, msg domain.NotificationMessage) error {
	ctx, cancel := context.WithTimeout(ctx, i.cfg.CommandTimeout)
	defer cancel()
	db := i.cli.Database(i.cfg.Database)
	coll := db.Collection(NewUserCollectionName(msg.UserUUID()))
	doc := documents.NotificationWriteModel{
		UUID:         msg.UUID(),
		UserUUID:     msg.UserUUID(),
		TrainingUUID: msg.TrainingUUID(),
		Trainer:      msg.Trainer(),
		Title:        msg.Title(),
		Content:      msg.Content(),
		Date:         msg.Date(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func NewInsertNotificationHandler(cli *mongo.Client, cfg InsertNotificationHandlerConfig) *InsertNotificationHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := InsertNotificationHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
