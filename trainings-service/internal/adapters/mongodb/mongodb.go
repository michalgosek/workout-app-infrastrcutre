package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
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
	trainer  *TrainerCommandHandler
	customer *CustomerCommandHandler
}

type QueryHandlers struct {
	trainer  *TrainerQueryHandler
	customer *CustomerQueryHandler
}

type MongoDB struct {
	cfg      Config
	commands CommandHandlers
	queries  QueryHandlers
}

func (m *MongoDB) UpsertWorkoutGroup(ctx context.Context, workoutGroup trainer.WorkoutGroup) error {
	return m.commands.trainer.UpsertWorkoutGroup(ctx, workoutGroup)
}
func (m *MongoDB) QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error) {
	return m.queries.trainer.QueryWorkoutGroup(ctx, groupUUID)
}

func (m *MongoDB) QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	return m.queries.trainer.QueryWorkoutGroups(ctx, trainerUUID)
}

func (m *MongoDB) DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error {
	return m.commands.trainer.DeleteWorkoutGroups(ctx, trainerUUID)
}

func (m *MongoDB) DeleteWorkoutGroup(ctx context.Context, groupUUID string) error {
	return m.commands.trainer.DeleteWorkoutGroup(ctx, groupUUID)
}

func (m *MongoDB) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	return m.commands.customer.UpsertCustomerWorkoutDay(ctx, workout)
}

func (m *MongoDB) QueryCustomerWorkoutDay(ctx context.Context, customerUUID string, trainerWorkoutGroupUUID string) (customer.WorkoutDay, error) {
	return m.queries.customer.QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkoutGroupUUID)
}

func (m *MongoDB) QueryCustomerWorkoutDays(ctx context.Context, customerUUID string) ([]customer.WorkoutDay, error) {
	return m.queries.customer.QueryCustomerWorkoutDays(ctx, customerUUID)
}

func (m *MongoDB) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, customerWorkoutDayUUID string) error {
	return m.commands.customer.DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkoutDayUUID)
}

func (m *MongoDB) DeleteCustomerWorkoutDays(ctx context.Context, customerUUID string) error {
	return m.commands.customer.DeleteCustomerWorkoutDays(ctx, customerUUID)
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
			customer: NewCustomerCommandHandler(mongoCLI, CustomerCommandHandlerConfig{
				Collection:     cfg.CustomerCollection,
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
			customer: NewCustomerQueryHandler(mongoCLI, CustomerQueryHandlerConfig{
				Collection:   cfg.CustomerCollection,
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

type document interface {
	CustomerWorkoutDocument | TrainerWorkoutGroupDocument
}

type filter interface {
	bson.M | bson.D
}

func updateOne[d document, f filter](ctx context.Context, coll *mongo.Collection, filter f, doc d) error {
	if coll == nil {
		return errors.New("empty collection")
	}
	update := bson.M{"$set": doc}
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}
