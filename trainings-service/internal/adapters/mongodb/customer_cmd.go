package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerCommandHandlerConfig struct {
	Collection     string
	Database       string
	Format         string
	CommandTimeout time.Duration
}

type CustomerCommandHandler struct {
	cfg CustomerCommandHandlerConfig
	cli *mongo.Client
}

func NewCustomerCommandHandler(cli *mongo.Client, cfg CustomerCommandHandlerConfig) *CustomerCommandHandler {
	t := CustomerCommandHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (c *CustomerCommandHandler) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	return coll.Drop(ctx)
}

func (c *CustomerCommandHandler) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	doc := CustomerWorkoutDocument{
		UUID:                    workout.UUID(),
		CustomerUUID:            workout.CustomerUUID(),
		TrainerWorkoutGroupUUID: workout.TrainerWorkoutGroupUUID(),
		Date:                    workout.Date().Format(c.cfg.Format),
	}
	filter := bson.M{"_id": workout.UUID()}
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)

	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	err := updateOne(ctx, coll, filter, doc)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (c *CustomerCommandHandler) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, customerWorkoutDayUUID string) error {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"_id": customerWorkoutDayUUID, "customer_uuid": customerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}
