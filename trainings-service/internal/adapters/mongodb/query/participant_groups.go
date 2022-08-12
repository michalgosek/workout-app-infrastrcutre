package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParticipantGroupsHandler struct {
	cli *mongo.Client
	cfg Config
}

func (t *ParticipantGroupsHandler) ParticipantGroups(ctx context.Context, participantUUID string) ([]query.ParticipantGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, t.cfg.QueryTimeout)
	defer cancel()

	f := bson.M{"participants._id": participantUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dd []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}
	gg := ConvertToParticipantGroups(dd...)
	return gg, nil
}

func NewParticipantGroupsHandler(cli *mongo.Client, cfg Config) *ParticipantGroupsHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := ParticipantGroupsHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
