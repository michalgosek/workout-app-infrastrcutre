package command

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateTrainerGroupHandler struct {
	cfg Config
	cli *mongo.Client
}

func (u *UpdateTrainerGroupHandler) Do(ctx context.Context, g *trainings.TrainingGroup) error {
	db := u.cli.Database(u.cfg.Database)
	coll := db.Collection(u.cfg.Collection)
	doc := documents.TrainingGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: documents.TrainerWriteModel{
			UUID: g.Trainer().UUID(),
			Name: g.Trainer().Name(),
		},
		Limit:        g.Limit(),
		Participants: UnmarshalToWriteModelParticipants(g.Participants()...),
	}
	filter := bson.M{"_id": g.UUID(), "trainer._id": g.Trainer().UUID()}
	update := bson.M{"$set": doc}
	ctx, cancel := context.WithTimeout(ctx, u.cfg.CommandTimeout)
	defer cancel()

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func NewUpdateTrainerGroupHandler(cli *mongo.Client, cfg Config) *UpdateTrainerGroupHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := UpdateTrainerGroupHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
