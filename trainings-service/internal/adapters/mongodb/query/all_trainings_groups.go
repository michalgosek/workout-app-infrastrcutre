package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AllTrainingsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *AllTrainingsHandler) Do(ctx context.Context) ([]query.TrainingGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, t.cfg.QueryTimeout)
	defer cancel()

	f := bson.D{}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dd []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}
	gg := UnmarshalToQueryTrainingGroups(dd...)
	return gg, nil
}

func NewAllTrainingsHandler(cli *mongo.Client, cfg Config) *AllTrainingsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := AllTrainingsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
