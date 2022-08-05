package query

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Config struct {
	Database     string
	Collection   string
	QueryTimeout time.Duration
}

type TrainerGroupHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *TrainerGroupHandler) Do(ctx context.Context, trainingUUID, trainerUUID string) (query.TrainerGroup, error) {
	f := bson.M{"_id": trainingUUID, "trainer._id": trainerUUID}
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return query.TrainerGroup{}, res.Err()
	}

	var doc documents.TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return query.TrainerGroup{}, err
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return query.TrainerGroup{}, nil
	}
	if err != nil {
		return query.TrainerGroup{}, err
	}
	m := UnmarshalToQueryTrainerWorkoutGroup(doc)
	return m, nil
}

func NewTrainerGroupHandler(cli *mongo.Client, cfg Config) *TrainerGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := TrainerGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
