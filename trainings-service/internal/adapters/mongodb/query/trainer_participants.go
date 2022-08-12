package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/command"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainerParticipantsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *TrainerParticipantsHandler) TrainerParticipants(ctx context.Context, trainerUUID string) ([]command.TrainerParticipant, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"trainer._id": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var docs []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, err
	}

	var out []command.TrainerParticipant
	for _, d := range docs {
		for _, p := range d.Participants {
			out = append(out, command.TrainerParticipant{
				TrainingUUID: d.UUID,
				TrainingName: d.Name,
				Trainer:      d.Trainer.Name,
				UserUUID:     p.UUID,
				Date:         d.Date,
			})
		}
	}
	return out, nil
}

func NewTrainerParticipantsHandler(cli *mongo.Client, cfg Config) *TrainerParticipantsHandler {
	if cli == nil {
		panic("nil mongo cli")
	}
	h := TrainerParticipantsHandler{
		cli: cli,
		cfg: cfg,
	}
	return &h
}
