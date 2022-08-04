package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertTrainerGroupHandler struct {
	cfg Config
	cli *mongo.Client
}

func (i *InsertTrainerGroupHandler) Do(ctx context.Context, g *trainings.TrainingGroup) error {
	ctx, cancel := context.WithTimeout(ctx, i.cfg.CommandTimeout)
	defer cancel()
	db := i.cli.Database(i.cfg.Database)
	coll := db.Collection(i.cfg.Collection)
	doc := documents.TrainingGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: documents.TrainerWriteModel{
			UUID: g.Trainer().UUID(),
			Name: g.Trainer().Name(),
		},
		Participants: UnmarshalToWriteModelParticipants(g.Participants()...),
		Limit:        g.Limit(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func NewInsertTrainerGroupHandler(cli *mongo.Client, cfg Config) *InsertTrainerGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := InsertTrainerGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
