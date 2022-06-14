package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerQueryHandlerConfig struct {
	Collection   string
	Database     string
	Format       string
	QueryTimeout time.Duration
}

type CustomerQueryHandler struct {
	cli *mongo.Client
	cfg CustomerQueryHandlerConfig
}

func NewCustomerQueryHandler(cli *mongo.Client, cfg CustomerQueryHandlerConfig) *CustomerQueryHandler {
	t := CustomerQueryHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (c *CustomerQueryHandler) QuerySchedule(ctx context.Context, customerUUID string) (customer.CustomerSchedule, error) {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"customer_uuid": customerUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return customer.CustomerSchedule{}, nil
	}
	if err != nil {
		return customer.CustomerSchedule{}, fmt.Errorf("find one failed: %v", err)
	}

	var dst CustomerScheduleDocument
	err = res.Decode(&dst)
	if err != nil {
		return customer.CustomerSchedule{}, fmt.Errorf("decoding failed: %v", err)
	}
	out := customer.UnmarshalFromDatabase(dst.UUID, dst.CustomerUUID, dst.Limit, dst.ScheduleUUIDs)
	return out, nil
}

func (c *CustomerQueryHandler) QuerySchedules(ctx context.Context, customerUUID string) ([]customer.CustomerSchedule, error) {
	return nil, nil
}
