package query

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-service/internal/adapters/mongodb/documents"
	"notification-service/internal/application/query"
	"time"
)

func NewUserCollectionName(UUID string) string {
	return fmt.Sprintf("%s_user", UUID)
}

type AllNotificationHandlerConfig struct {
	Database       string
	CommandTimeout time.Duration
}

type AllNotificationHandler struct {
	cfg AllNotificationHandlerConfig
	cli *mongo.Client
}

func (a *AllNotificationHandler) AllNotifications(ctx context.Context, UUID string) ([]query.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, a.cfg.CommandTimeout)
	defer cancel()

	db := a.cli.Database(a.cfg.Database)
	name := NewUserCollectionName(UUID)
	coll := db.Collection(name)
	f := bson.M{"user_uuid": UUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var docs []documents.NotificationWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	out := ConvertToQueryAllNotifications(docs...)
	return out, nil
}

func ConvertToQueryAllNotifications(dd ...documents.NotificationWriteModel) []query.Notification {
	var out []query.Notification
	for _, d := range dd {
		out = append(out, query.Notification{
			UUID:         d.UUID,
			UserUUID:     d.UserUUID,
			TrainingUUID: d.TrainingUUID,
			Title:        d.Title,
			Trainer:      d.Trainer,
			Content:      d.Content,
			Date:         d.Date.Format(query.UIFormat),
		})
	}
	return out
}

func NewAllNotificationHandler(cli *mongo.Client, cfg AllNotificationHandlerConfig) *AllNotificationHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := AllNotificationHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
