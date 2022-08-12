package query

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DuplicateTrainingGroupHandler struct {
	cli *mongo.Client
	cfg Config
}

func (d DuplicateTrainingGroupHandler) IsTrainingGroupExists(ctx context.Context, g *trainings.TrainingGroup) (bool, error) {
	f := bson.M{"trainer._id": g.Trainer().UUID(), "description": g.Description(), "name": g.Name(), "date": g.Date()}
	db := d.cli.Database(d.cfg.Database)
	coll := db.Collection(d.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), d.cfg.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewDuplicateTrainingGroupHandler(cli *mongo.Client, cfg Config) *DuplicateTrainingGroupHandler {
	if cli == nil {
		panic("nil mongo cli")
	}
	h := DuplicateTrainingGroupHandler{
		cli: cli,
		cfg: cfg,
	}
	return &h
}
