package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainingsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *TrainingsHandler) Do(ctx context.Context, trainerUUID string) ([]query.TrainerWorkoutGroup, error) {
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
	gg := UnmarshalToQueryTrainerWorkoutGroups(dd...)
	return gg, nil
}

func NewTrainingsHandler(cli *mongo.Client, cfg Config) *TrainingsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := TrainingsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}

func UnmarshalToQueryTrainerWorkoutGroups(dd ...documents.TrainingGroupWriteModel) []query.TrainerWorkoutGroup {
	var out []query.TrainerWorkoutGroup
	for _, d := range dd {
		g := UnmarshalToQueryTrainerWorkoutGroup(d)
		out = append(out, g)
	}
	return out
}
