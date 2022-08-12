package query

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainingGroupHandler struct {
	cfg Config
	cli *mongo.Client
}

func (t TrainingGroupHandler) TrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)

	f := bson.M{"_id": trainingUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainings.TrainingGroup{}, nil
	}
	if err != nil {
		return trainings.TrainingGroup{}, err
	}

	var doc documents.TrainingGroupWriteModel
	err = res.Decode(&doc)
	if err != nil {
		return trainings.TrainingGroup{}, err
	}

	var pp []trainings.DatabaseTrainingGroupParticipant
	for _, p := range doc.Participants {
		pp = append(pp, trainings.DatabaseTrainingGroupParticipant{UUID: p.UUID, Name: p.Name})
	}
	g := trainings.ConvertTrainingGroupFromDatabase(trainings.DatabaseTrainingGroup{
		UUID:        doc.UUID,
		Name:        doc.Name,
		Description: doc.Description,
		Limit:       doc.Limit,
		Date:        doc.Date,
		Trainer: trainings.DatabaseTrainingGroupTrainer{
			UUID: doc.Trainer.UUID,
			Name: doc.Trainer.Name,
		},
		Participants: pp,
	})
	return g, nil
}

func NewTrainingGroupHandler(cli *mongo.Client, cfg Config) *TrainingGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := TrainingGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
