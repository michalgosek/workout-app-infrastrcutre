package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainerCommandHandlerConfig struct {
	Collection     string
	Database       string
	Format         string
	CommandTimeout time.Duration
}

type TrainerCommandHandler struct {
	cli *mongo.Client
	cfg TrainerCommandHandlerConfig
}

func NewTrainerCommandHandler(cli *mongo.Client, cfg TrainerCommandHandlerConfig) *TrainerCommandHandler {
	t := TrainerCommandHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (m *TrainerCommandHandler) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	return coll.Drop(ctx)
}

func (m *TrainerCommandHandler) UpsertWorkoutGroup(ctx context.Context, schedule trainer.WorkoutGroup) error {
	doc := TrainerWorkoutGroupDocument{
		UUID:          schedule.UUID(),
		TrainerUUID:   schedule.TrainerUUID(),
		Limit:         schedule.Limit(),
		CustomerUUIDs: schedule.CustomerUUIDs(),
		Name:          schedule.Name(),
		Desc:          schedule.Desc(),
		Date:          schedule.Date().Format(m.cfg.Format),
	}
	f := bson.M{"_id": schedule.UUID(), "trainer_uuid": schedule.TrainerUUID()}
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	err := updateOne(ctx, coll, f, doc)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (m *TrainerCommandHandler) DeleteWorkoutGroups(ctx context.Context, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return fmt.Errorf("delete many failed: %v", err)
	}
	return nil
}

func (m *TrainerCommandHandler) DeleteWorkoutGroup(ctx context.Context, groupUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": groupUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}
