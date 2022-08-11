package testutil

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"notification-service/internal/adapters/mongodb/documents"
	"notification-service/internal/domain"
	"time"
)

func NewTestUserNotification(userUUID, trainingUUID string) domain.NotificationMessage {
	n, err := domain.NewNotificationMessage(domain.NotificationMessageData{
		UserUUID:     userUUID,
		TrainingUUID: trainingUUID,
		Title:        fmt.Sprintf("title_for_user_%s", userUUID),
		Content:      fmt.Sprintf("content_for_user_%s", userUUID),
		Date:         NewTestStaticTime(),
		Trainer:      "John Doe",
	})
	if err != nil {
		panic(err)
	}
	return n
}

func NewTestStaticTime() time.Time {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	return ts
}

type TestClientConfig struct {
	Addr           string
	CommandTimeout time.Duration
	QueryTimeout   time.Duration
	Database       string
	Collection     string
}

type TestClient struct {
	cli *mongo.Client
	cfg TestClientConfig
}

func NewTestMongoClient(cfg TestClientConfig) *TestClient {
	opts := options.Client()
	opts.ApplyURI(cfg.Addr)
	opts.SetConnectTimeout(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		panic(err)
	}
	err = cli.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	tc := TestClient{
		cli: cli,
		cfg: cfg,
	}
	return &tc
}

func (t *TestClient) Disconnect(ctx context.Context) error {
	err := t.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestClient) Drop(ctx context.Context) error {
	db := t.cli.Database(t.cfg.Database)
	err := db.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TestClient) FindInsertedNotification(UUID string) (documents.NotificationWriteModel, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.QueryTimeout)
	defer cancel()
	f := bson.M{"user_uuid": UUID}
	sr := coll.FindOne(ctx, f)
	if sr.Err() != nil {
		return documents.NotificationWriteModel{}, sr.Err()
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc documents.NotificationWriteModel
	err := sr.Decode(&doc)
	if err != nil {
		return documents.NotificationWriteModel{}, err
	}
	return doc, nil
}

func (t *TestClient) FindInsertedNotifications(UUID string) ([]documents.NotificationWriteModel, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.QueryTimeout)
	defer cancel()
	f := bson.M{"user_uuid": UUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var docs []documents.NotificationWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
