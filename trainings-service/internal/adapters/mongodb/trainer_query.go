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

func (t *TrainerQueryHandler) QueryWorkoutGroup(ctx context.Context, UUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"_id": UUID, "trainer_uuid": trainerUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainer.WorkoutGroup{}, nil
	}
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("find one failed: %v", err)
	}

	var doc TrainerWorkoutGroupDocument
	err = res.Decode(&doc)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("decoding failed: %v", err)
	}
	date, err := time.Parse(t.cfg.Format, doc.Date)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	out, err := trainer.UnmarshalFromDatabase(doc.UUID, doc.TrainerUUID, doc.Name, doc.Desc, doc.CustomerUUIDs, date, doc.Limit)
	return out, nil
}

func (t *TrainerQueryHandler) QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("find failed: %v", err)
	}

	var docs []TrainerWorkoutGroupDocument
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %v", err)
	}

	workouts, err := convertToDomainWorkoutGroups(t.cfg.Format, docs...)
	if err != nil {
		return nil, fmt.Errorf("converting docs to domain workout groups failed: %v", err)
	}
	return workouts, nil
}

func convertToDomainWorkoutGroups(format string, docs ...TrainerWorkoutGroupDocument) ([]trainer.WorkoutGroup, error) {
	var workouts []trainer.WorkoutGroup

	for _, d := range docs {
		date, err := time.Parse(format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		workout, err := trainer.UnmarshalFromDatabase(d.UUID, d.TrainerUUID, d.Name, d.Desc, d.CustomerUUIDs, date, d.Limit)
		if err != nil {
			return nil, fmt.Errorf("unmarshal from database failed: %v", err)
		}
		workouts = append(workouts, workout)
	}
	return workouts, nil
}
