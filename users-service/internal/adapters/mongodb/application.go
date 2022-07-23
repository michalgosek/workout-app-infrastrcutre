package mongodb

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/query"
	rm "github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Commands struct {
	*command.InsertUserHandler
}

type Queries struct {
	*query.UserHandler
}

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
	cfg      Config
	cli      *mongo.Client
	commands *Commands
	queries  *Queries
}

func (r *Repository) findOne(ctx context.Context, f bson.M) (documents.UserWriteModel, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
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

func (r *Repository) InsertUser(ctx context.Context, user *domain.User) error {
	return r.commands.Do(ctx, user)
}

func (r *Repository) User(ctx context.Context, UUID string) (rm.User, error) {
	return r.queries.User(ctx, UUID)
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

	var dst documents.UserWriteModel
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

func (r *Repository) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.TODO(), r.cfg.Timeouts.ConnectionTimeout)
	defer cancel()
	err := r.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(cfg Config) (*Repository, error) {
	cli, err := NewClient(cfg.Addr, cfg.Timeouts.ConnectionTimeout)
	if err != nil {
		return nil, err
	}
	commandCfg := command.Config{
		Database:       cfg.Database,
		Collection:     cfg.Collection,
		CommandTimeout: cfg.Timeouts.CommandTimeout,
	}
	queryCfg := query.Config{
		Database:     cfg.Database,
		Collection:   cfg.Collection,
		QueryTimeout: cfg.Timeouts.QueryTimeout,
	}
	r := Repository{
		cfg: cfg,
		cli: cli,
		commands: &Commands{
			InsertUserHandler: command.NewInsertUserHandler(cli, commandCfg),
		},
		queries: &Queries{
			UserHandler: query.NewUserHandler(cli, queryCfg),
		},
	}
	return &r, nil
}

func NewClient(addr string, d time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(d)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	err = cli.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return cli, nil
}
