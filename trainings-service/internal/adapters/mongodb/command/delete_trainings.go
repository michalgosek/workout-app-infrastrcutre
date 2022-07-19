package command

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteTrainingsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (d *DeleteTrainingsHandler) Do(ctx context.Context, trainerUUID string) error {
	f := bson.M{"trainer._id": trainerUUID}
	ctx, cancel := context.WithTimeout(ctx, d.cfg.CommandTimeout)
	defer cancel()
	db := d.cli.Database(d.cfg.Database)
	coll := db.Collection(d.cfg.Collection)

	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return nil
	}
	return nil
}

func NewDeleteTrainingsHandler(cli *mongo.Client, cfg Config) *DeleteTrainingsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := DeleteTrainingsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
