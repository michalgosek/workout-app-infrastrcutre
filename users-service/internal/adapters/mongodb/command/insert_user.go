package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Config struct {
	Database       string
	Collection     string
	CommandTimeout time.Duration
}

type InsertUserHandler struct {
	cli *mongo.Client
	cfg Config
}

func (i *InsertUserHandler) Do(ctx context.Context, u *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, i.cfg.CommandTimeout)
	defer cancel()
	db := i.cli.Database(i.cfg.Database)
	coll := db.Collection(i.cfg.Collection)
	doc := UnmarshalToUserWriteModel(*u)
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func NewInsertUserHandler(cli *mongo.Client, cfg Config) *InsertUserHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := InsertUserHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
