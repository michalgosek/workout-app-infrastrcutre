package query

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Config struct {
	Database     string
	Collection   string
	QueryTimeout time.Duration
}

type UserHandler struct {
	cli *mongo.Client
	cfg Config
}

func (u *UserHandler) findOne(ctx context.Context, f bson.M) (documents.UserWriteModel, error) {
	db := u.cli.Database(u.cfg.Database)
	coll := db.Collection(u.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), u.cfg.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return documents.UserWriteModel{}, res.Err()
	}

	var dst documents.UserWriteModel
	err := res.Decode(&dst)
	if err != nil {
		return documents.UserWriteModel{}, nil
	}
	return dst, nil
}

func (u *UserHandler) User(ctx context.Context, UUID string) (query.User, error) {
	f := bson.M{"_id": UUID}
	doc, err := u.findOne(ctx, f)
	if err != nil {
		return query.User{}, nil
	}
	g := query.User{
		UUID:  doc.UUID,
		Name:  doc.Name,
		Role:  doc.Role,
		Email: doc.Email,
	}
	return g, nil
}

func NewUserHandler(cli *mongo.Client, cfg Config) *UserHandler {
	if cli == nil {
		panic("nil mongo client")
	}
	h := UserHandler{
		cfg: cfg,
		cli: cli,
	}
	return &h
}
