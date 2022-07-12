package customer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueryHandlerConfig struct {
	Collection   string
	Database     string
	Format       string
	QueryTimeout time.Duration
}

type QueryHandler struct {
	cli *mongo.Client
	cfg QueryHandlerConfig
}

func NewQueryHandler(cli *mongo.Client, cfg QueryHandlerConfig) *QueryHandler {
	t := QueryHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (c *QueryHandler) QueryCustomerWorkoutDay(ctx context.Context, customerUUID, GroupUUID string) (customer.WorkoutDay, error) {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"customer_uuid": customerUUID, "trainer_workout_group_uuid": GroupUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return customer.WorkoutDay{}, nil
	}
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("find one failed: %v", err)
	}

	var dst WorkoutDayDocument
	err = res.Decode(&dst)
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("decoding failed: %v", err)
	}

	date, err := time.Parse(c.cfg.Format, dst.Date)
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	out, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
		WorkoutDayUUID: dst.UUID,
		TrainerUUID:    dst.TrainerUUID,
		GroupUUID:      dst.GroupUUID,
		CustomerUUID:   dst.CustomerUUID,
		CustomerName:   dst.CustomerName,
		Date:           date,
	})
	if err != nil {
		return customer.WorkoutDay{}, fmt.Errorf("unmarshal failed: %v", err)
	}
	return out, nil
}

func (c *QueryHandler) QueryCustomerWorkoutDays(ctx context.Context, customerUUID string) ([]customer.WorkoutDay, error) {
	db := c.cli.Database(c.cfg.Database)
	coll := db.Collection(c.cfg.Collection)
	f := bson.M{"customer_uuid": customerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("find failed: %v", err)
	}

	var docs []WorkoutDayDocument
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}

	days, err := convertDocumentsToWorkoutDays(c.cfg.Format, docs...)
	return days, nil
}

func convertDocumentsToWorkoutDays(format string, docs ...WorkoutDayDocument) ([]customer.WorkoutDay, error) {
	var days []customer.WorkoutDay
	for _, d := range docs {
		date, err := time.Parse(format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		day, err := customer.UnmarshalFromDatabase(customer.UnmarshalFromDatabaseArgs{
			WorkoutDayUUID: d.UUID,
			TrainerUUID:    d.TrainerUUID,
			GroupUUID:      d.GroupUUID,
			CustomerUUID:   d.CustomerUUID,
			CustomerName:   d.CustomerName,
			Date:           date,
		})
		if err != nil {
			return nil, fmt.Errorf("unmarshal from database failed: %v", err)
		}
		days = append(days, day)
	}
	return days, nil
}
