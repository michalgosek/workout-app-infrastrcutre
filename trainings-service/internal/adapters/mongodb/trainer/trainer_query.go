package trainer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
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

func (t *QueryHandler) QueryTrainerWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"_id": groupUUID}
	res := coll.FindOne(ctx, f)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainer.WorkoutGroup{}, nil
	}
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("find one failed: %v", err)
	}

	var doc WorkoutGroupDocument
	err = res.Decode(&doc)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("decoding failed: %v", err)
	}
	date, err := time.Parse(t.cfg.Format, doc.Date)
	if err != nil {
		return trainer.WorkoutGroup{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	group, err := trainer.UnmarshalFromDatabase(
		doc.UUID,
		doc.TrainerUUID,
		doc.TrainerName,
		doc.WorkoutName,
		doc.WorkoutDesc,
		doc.CustomerUUIDs,
		date,
		doc.Limit,
	)
	return group, nil
}

func (t *QueryHandler) QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("find failed: %v", err)
	}

	var docs []WorkoutGroupDocument
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %v", err)
	}

	groups, err := convertDocumentsToWorkoutGroups(t.cfg.Format, docs...)
	if err != nil {
		return nil, fmt.Errorf("converting docs to domain workout groups failed: %v", err)
	}
	return groups, nil
}

func convertDocumentsToWorkoutGroups(format string, docs ...WorkoutGroupDocument) ([]trainer.WorkoutGroup, error) {
	var groups []trainer.WorkoutGroup
	for _, d := range docs {
		date, err := time.Parse(format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		group, err := trainer.UnmarshalFromDatabase(
			d.UUID,
			d.TrainerUUID,
			d.TrainerName,
			d.WorkoutName,
			d.WorkoutDesc,
			d.CustomerUUIDs,
			date,
			d.Limit,
		)
		if err != nil {
			return nil, fmt.Errorf("unmarshal from database failed: %v", err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}
