package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertTrainingGroupHandler struct {
	cfg Config
	cli *mongo.Client
}

func (i *InsertTrainingGroupHandler) InsertTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
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
		Participants: ConvertToWriteModelParticipants(g.Participants()...),
		Limit:        g.Limit(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func NewInsertTrainingGroupHandler(cli *mongo.Client, cfg Config) *InsertTrainingGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := InsertTrainingGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
