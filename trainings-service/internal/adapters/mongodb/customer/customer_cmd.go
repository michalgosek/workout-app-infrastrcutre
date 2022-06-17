package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommandHandlerConfig struct {
	Collection     string
	Database       string
	Format         string
	CommandTimeout time.Duration
}

type CommandHandler struct {
	cfg CommandHandlerConfig
	cli *mongo.Client
}

func NewCommandHandler(cli *mongo.Client, cfg CommandHandlerConfig) *CommandHandler {
	t := CommandHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (c *CommandHandler) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	return coll.Drop(ctx)
}

func (c *CommandHandler) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	doc := WorkoutDocument{
		UUID:                    workout.UUID(),
		CustomerUUID:            workout.CustomerUUID(),
		TrainerWorkoutGroupUUID: workout.GroupUUID(),
		Date:                    workout.Date().Format(c.cfg.Format),
	}

	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()

	f := bson.M{"_id": workout.UUID()}
	update := bson.M{"$set": doc}
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := coll.UpdateOne(ctx, f, update, opts)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (c *CommandHandler) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, workoutDayUUID string) error {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"_id": workoutDayUUID, "customer_uuid": customerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}

func (c *CommandHandler) DeleteCustomerWorkoutDays(ctx context.Context, customerUUID string) error {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"customer_uuid": customerUUID}
	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return fmt.Errorf("delete many failed: %v", err)
	}
	return nil
}
