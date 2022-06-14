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

func (c *CustomerQueryHandler) QueryCustomerWorkoutDay(ctx context.Context, customerUUID, trainerWorkoutGroupUUID string) (customer.WorkoutDay, error) {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"customer_uuid": customerUUID, "trainer_workout_group_uuid": trainerWorkoutGroupUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return customer.WorkoutDay{}, nil
	}
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("find one failed: %v", err)
	}

	var dst CustomerWorkoutDocument
	err = res.Decode(&dst)
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("decoding failed: %v", err)
	}

	date, err := time.Parse(c.cfg.Format, dst.Date)
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	out, err := customer.UnmarshalFromDatabase(dst.UUID, dst.TrainerWorkoutGroupUUID, dst.CustomerUUID, date)
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("unmarshal failed: %v", err)
	}
	return out, nil
}
