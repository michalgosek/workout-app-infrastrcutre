package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainerQueryHandlerConfig struct {
	Collection   string
	Database     string
	Format       string
	QueryTimeout time.Duration
}

type TrainerQueryHandler struct {
	cli *mongo.Client
	cfg TrainerQueryHandlerConfig
}

func NewTrainerQueryHandler(cli *mongo.Client, cfg TrainerQueryHandlerConfig) *TrainerQueryHandler {
	t := TrainerQueryHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (t *TrainerQueryHandler) QueryTrainerSchedule(ctx context.Context, UUID, trainerUUID string) (trainer.TrainerSchedule, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"_id": UUID, "trainer_uuid": trainerUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainer.TrainerSchedule{}, nil
	}
	if err != nil {
		return trainer.TrainerSchedule{}, fmt.Errorf("find one failed: %v", err)
	}

	var dst TrainerScheduleDocument
	err = res.Decode(&dst)
	if err != nil {
		return trainer.TrainerSchedule{}, fmt.Errorf("decoding failed: %v", err)
	}
	date, err := time.Parse(t.cfg.Format, dst.Date)
	if err != nil {
		return trainer.TrainerSchedule{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	out := trainer.UnmarshalFromDatabase(dst.UUID, dst.TrainerUUID, dst.Name, dst.Desc, dst.CustomerUUIDs, date, dst.Limit)
	return out, nil
}

func (t *TrainerQueryHandler) QueryTrainerSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("find failed: %v", err)
	}

	var dst []TrainerScheduleDocument
	err = cur.All(ctx, &dst)
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %v", err)
	}
	var out []trainer.TrainerSchedule
	for _, d := range dst {
		date, err := time.Parse(t.cfg.Format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		s := trainer.UnmarshalFromDatabase(d.UUID, d.TrainerUUID, d.Name, d.Desc, d.CustomerUUIDs, date, d.Limit)
		out = append(out, s)
	}
	return out, nil
}
