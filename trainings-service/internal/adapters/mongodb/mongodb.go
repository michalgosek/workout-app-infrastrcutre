package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Addr               string
	Database           string
	TrainerCollection  string
	CustomerCollection string
	CommandTimeout     time.Duration
	QueryTimeout       time.Duration
	ConnectionTimeout  time.Duration
	Format             string
}

type CommandHandlers struct {
	trainer *TrainerCommandHandler
}

type QueryHandlers struct {
	trainer *TrainerQueryHandler
}

type MongoDB struct {
	cfg      Config
	commands CommandHandlers
	queries  QueryHandlers
}

func (m *MongoDB) UpsertTrainerSchedule(ctx context.Context, schedule trainer.TrainerSchedule) error {
	return m.commands.trainer.UpsertSchedule(ctx, schedule)
}
func (m *MongoDB) QueryTrainerSchedule(ctx context.Context, UUID, trainerUUID string) (trainer.TrainerSchedule, error) {
	return m.queries.trainer.QueryTrainerSchedule(ctx, UUID, trainerUUID)
}

func (m *MongoDB) QueryTrainerSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	return m.queries.trainer.QueryTrainerSchedules(ctx, trainerUUID)
}

func (m *MongoDB) CancelTrainerSchedules(ctx context.Context, trainerUUID string) error {
	return m.commands.trainer.CancelSchedules(ctx, trainerUUID)
}

func (m *MongoDB) CancelTrainerSchedule(ctx context.Context, UUID, trainerUUID string) error {
	return m.commands.trainer.CancelSchedule(ctx, UUID, trainerUUID)
}

func NewMongoDB(cfg Config) (*MongoDB, error) {
	mongoCLI, err := newMongoClient(cfg.Addr, cfg.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("creating mongo cli failed: %v", err)
	}
	m := MongoDB{
		cfg: cfg,
		commands: CommandHandlers{
			trainer: NewTrainerCommandHandler(mongoCLI, TrainerCommandHandlerConfig{
				Collection:     cfg.TrainerCollection,
				Database:       cfg.Database,
				Format:         cfg.Format,
				CommandTimeout: cfg.CommandTimeout,
			}),
		},
		queries: QueryHandlers{
			trainer: NewTrainerQueryHandler(mongoCLI, TrainerQueryHandlerConfig{
				Collection:   cfg.TrainerCollection,
				Database:     cfg.Database,
				Format:       cfg.Format,
				QueryTimeout: cfg.CommandTimeout,
			}),
		},
	}
	return &m, nil
}

func newMongoClient(addr string, timeout time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	mongoCLI, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("mongo client creation failed: %v", err)
	}
	err = mongoCLI.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client connection failed: %v", err)
	}
	err = mongoCLI.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("mongo client ping req failed: %v", err)
	}
	return mongoCLI, nil
}
