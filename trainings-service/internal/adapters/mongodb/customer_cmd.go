package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *CustomerCommandHandler) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	return coll.Drop(ctx)
}

func (c *CustomerCommandHandler) UpsertSchedule(ctx context.Context, schedule customer.CustomerSchedule) error {
	doc := CustomerScheduleDocument{
		UUID:          schedule.UUID(),
		CustomerUUID:  schedule.CustomerUUID(),
		Limit:         schedule.Limit(),
		ScheduleUUIDs: schedule.ScheduleUUIDs(),
	}
	filter := bson.M{"_id": schedule.UUID()}
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	update := bson.M{
		"$set": doc,
	}
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (c *CustomerCommandHandler) DeleteSchedule(ctx context.Context, customerUUID, scheduleUUID string) error {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"_id": scheduleUUID, "customer_uuid": customerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}

func (c *CustomerCommandHandler) CancelSchedules(ctx context.Context, customerUUID string) error {
	return nil
}
