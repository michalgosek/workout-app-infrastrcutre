package mongodb

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Timeouts struct {
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
}

type Config struct {
	Addr       string
	Database   string
	Collection string
	Timeouts   Timeouts
}

type Repository struct {
	cfg Config
	cli *mongo.Client
}

func (r *Repository) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.ConnectionTimeout)
	defer cancel()
	err := r.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertUser(ctx context.Context, u *domain.User) error {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.CommandTimeout)
	defer cancel()

	doc := UnmarshalToUserWriteModel(*u)
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return nil
	}
	return nil
}

func (r *Repository) findOne(ctx context.Context, f bson.M) (UserWriteModel, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return UserWriteModel{}, res.Err()
	}

	var dst UserWriteModel
	err := res.Decode(&dst)
	if err != nil {
		return UserWriteModel{}, nil
	}
	return dst, nil
}

func (r *Repository) QueryUserWithEmail(ctx context.Context, email string) (domain.User, error) {
	f := bson.M{"email": email}
	doc, err := r.findOne(ctx, f)
	if err != nil {
		return domain.User{}, nil
	}
	user := domain.UnmarshalUserFromDatabase(domain.DatabaseUser{
		UUID:           doc.UUID,
		Active:         doc.Active,
		Role:           doc.Role,
		Name:           doc.Name,
		Email:          doc.Email,
		LastActiveDate: doc.LastActiveDate,
	})
	return user, nil
}

func (r *Repository) User(ctx context.Context, UUID string) (query.User, error) {
	f := bson.M{"_id": UUID}
	doc, err := r.findOne(ctx, f)
	if err != nil {
		return query.User{}, nil
	}
	u := query.User{
		Name:  doc.Name,
		Role:  doc.Role,
		Email: doc.Email,
	}
	return u, nil
}

func (r *Repository) QueryUser(ctx context.Context, UUID string) (domain.User, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	f := bson.M{"_id": UUID}
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return domain.User{}, res.Err()
	}

	var dst UserWriteModel
	err := res.Decode(&dst)
	if err != nil {
		return domain.User{}, nil
	}

	u := domain.UnmarshalUserFromDatabase(domain.DatabaseUser{
		UUID:           dst.UUID,
		Active:         dst.Active,
		Role:           dst.Role,
		Name:           dst.Name,
		Email:          dst.Email,
		LastActiveDate: dst.LastActiveDate,
	})
	return u, nil
}

func NewRepository(cfg Config) (*Repository, error) {
	cli, err := NewClient(cfg.Addr, cfg.Timeouts.ConnectionTimeout)
	if err != nil {
		return nil, err
	}
	m := Repository{
		cli: cli,
		cfg: cfg,
	}
	return &m, nil
}
