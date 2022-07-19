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

type TrainingHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *TrainingHandler) Do(ctx context.Context, trainingUUID, trainerUUID string) (query.TrainerWorkoutGroup, error) {
	f := bson.M{"_id": trainingUUID, "trainer._id": trainerUUID}
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), t.cfg.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return query.TrainerWorkoutGroup{}, res.Err()
	}

	var doc documents.TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return query.TrainerWorkoutGroup{}, err
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return query.TrainerWorkoutGroup{}, nil
	}
	if err != nil {
		return query.TrainerWorkoutGroup{}, err
	}
	m := UnmarshalToQueryTrainerWorkoutGroup(doc)
	return m, nil
}

func NewTrainingHandler(cli *mongo.Client, cfg Config) *TrainingHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := TrainingHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
