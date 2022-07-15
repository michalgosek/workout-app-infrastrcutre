package trainer

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
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
	cli *mongo.Client
	cfg CommandHandlerConfig
}

func NewCommandHandler(cli *mongo.Client, cfg CommandHandlerConfig) *CommandHandler {
	t := CommandHandler{
		cli: cli,
		cfg: cfg,
	}
	return &t
}

func (m *CommandHandler) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	return coll.Drop(ctx)
}

func (m *CommandHandler) UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error {
	var details []WorkoutGroupParticipantWriteModel
	for _, d := range group.CustomerDetails() {
		details = append(details, WorkoutGroupParticipantWriteModel{UUID: d.UUID(), Name: d.Name()})
	}

	doc := WorkoutGroupWriteModel{
		UUID:         group.UUID(),
		TrainerUUID:  group.TrainerUUID(),
		TrainerName:  group.TrainerName(),
		Limit:        group.Limit(),
		Participants: details,
		Name:         group.Name(),
		Description:  group.Description(),
		Date:         group.Date().Format(m.cfg.Format),
	}

	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	update := bson.M{"$set": doc}
	opts := options.Update()
	opts.SetUpsert(true)
	f := bson.M{"_id": group.UUID(), "trainer_uuid": group.TrainerUUID()}
	_, err := coll.UpdateOne(ctx, f, update, opts)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (m *CommandHandler) DeleteTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": groupUUID, "trainer_uuid": trainerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}

func (m *CommandHandler) DeleteTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return fmt.Errorf("delete many failed: %v", err)
	}
	return nil
}
