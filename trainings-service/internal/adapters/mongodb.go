package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Filter struct {
	ScheduleUUID string
	ScheduleName string
	TrainerUUID  string
}

type MongoDBConfig struct {
	Addr              string
	Database          string
	Collection        string
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
}

type TrainerSchedulesMongoDB struct {
	cli *mongo.Client
	cfg MongoDBConfig
}

type TrainerScheduleDocument struct {
	UUID          string    `bson:"_id"`
	TrainerUUID   string    `bson:"trainer_uuid"`
	Limit         int       `bson:"limit"`
	CustomerUUIDs []string  `bson:"customer_uuids"`
	Name          string    `bson:"name"`
	Desc          string    `bson:"desc"`
	Date          time.Time `bson:"date"`
}

func (t *TrainerSchedulesMongoDB) Disconnect(ctx context.Context) error {
	return t.cli.Disconnect(ctx)
}

func (t *TrainerSchedulesMongoDB) Ping(ctx context.Context) error {
	err := t.cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("ping request failed: %v", err)
	}
	return err
}

func (t *TrainerSchedulesMongoDB) DropCollection(ctx context.Context) error {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	return coll.Drop(ctx)
}

func (t *TrainerSchedulesMongoDB) UpsertSchedule(ctx context.Context, schedule trainer.TrainerSchedule) error {
	db := t.cli.Database(t.cfg.Database)
	coll := db.Collection(t.cfg.Collection)
	opts := options.Update()
	opts.SetUpsert(true)
	doc := TrainerScheduleDocument{
		UUID:          schedule.UUID(),
		TrainerUUID:   schedule.TrainerUUID(),
		Limit:         schedule.Limit(),
		CustomerUUIDs: schedule.CustomerUUIDs(),
		Name:          schedule.Name(),
		Desc:          schedule.Desc(),
		Date:          schedule.Date(),
	}
	update := bson.M{
		"$set": doc,
	}
	filter := bson.M{"_id": schedule.UUID(), "trainer_uuid": schedule.TrainerUUID()}
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update document failed: %v", err)
	}
	return nil
}
func (t *TrainerSchedulesMongoDB) QuerySchedule(ctx context.Context, UUID, trainerUUID string) (trainer.TrainerSchedule, error) {
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
	out := trainer.UnmarshalFromDatabase(dst.UUID, dst.TrainerUUID, dst.Name, dst.Desc, dst.CustomerUUIDs, dst.Date, dst.Limit)
	return out, nil
}

func (m *TrainerSchedulesMongoDB) QuerySchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
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
		s := trainer.UnmarshalFromDatabase(d.UUID, d.TrainerUUID, d.Name, d.Desc, d.CustomerUUIDs, d.Date, d.Limit)
		out = append(out, s)
	}
	return out, nil
}

func (m *TrainerSchedulesMongoDB) CancelSchedules(ctx context.Context, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return fmt.Errorf("delete many failed: %v", err)
	}
	return nil
}

func (m *TrainerSchedulesMongoDB) CancelSchedule(ctx context.Context, UUID, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": UUID, "trainer_uuid": trainerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}

func NewTrainerSchedulesMongoDB(cfg MongoDBConfig) (*TrainerSchedulesMongoDB, error) {
	opts := options.Client()
	opts.SetConnectTimeout(cfg.ConnectionTimeout)
	opts.ApplyURI(cfg.Addr)

	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("new client failed: %v", err)
	}
	err = cli.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connect failed: %v", err)
	}

	m := TrainerSchedulesMongoDB{
		cli: cli,
		cfg: cfg,
	}
	return &m, nil
}
