package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainerGroupsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *TrainerGroupsHandler) TrainerGroups(ctx context.Context, trainerUUID string) ([]query.TrainerGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, t.cfg.QueryTimeout)
	defer cancel()

	f := bson.M{"trainer._id": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dd []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}
	gg := ConvertToTrainerWorkoutGroups(dd...)
	return gg, nil
}

func NewTrainerGroupsHandler(cli *mongo.Client, cfg Config) *TrainerGroupsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := TrainerGroupsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
