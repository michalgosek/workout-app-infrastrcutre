package command

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteTrainerGroupHandler struct {
	cli *mongo.Client
	cfg Config
}

func (d *DeleteTrainerGroupHandler) Do(ctx context.Context, trainingUUID, trainerUUID string) error {
	f := bson.M{"_id": trainingUUID, "trainer._id": trainerUUID}
	ctx, cancel := context.WithTimeout(ctx, d.cfg.CommandTimeout)
	defer cancel()
	db := d.cli.Database(d.cfg.Database)
	coll := db.Collection(d.cfg.Collection)

	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return nil
	}
	return nil
}

func NewDeleteTrainerGroupHandler(cli *mongo.Client, cfg Config) *DeleteTrainerGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := DeleteTrainerGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
