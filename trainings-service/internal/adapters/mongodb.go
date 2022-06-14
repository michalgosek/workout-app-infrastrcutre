package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type TrainerScheduleDocument struct {
	UUID          string   `bson:"_id"`
	TrainerUUID   string   `bson:"trainer_uuid"`
	Limit         int      `bson:"limit"`
	CustomerUUIDs []string `bson:"customer_uuids"`
	Name          string   `bson:"name"`
	Desc          string   `bson:"desc"`
	Date          string   `bson:"date"`
}

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
	Format            string
}

type MongoDB struct {
	cli *mongo.Client
	cfg MongoDBConfig
}

func (m *MongoDB) QueryCustomerSchedule(ctx context.Context, customerUUID string) (customer.CustomerSchedule, error) {
	return customer.CustomerSchedule{}, nil
}

func (m *MongoDB) QueryCustomerSchedules(ctx context.Context, customerUUID string) ([]customer.CustomerSchedule, error) {
	return []customer.CustomerSchedule{}, nil
}

func (m *MongoDB) UpsertCustomerSchedule(ctx context.Context, schedule customer.CustomerSchedule) error {
	return nil
}

func (m *MongoDB) UpsertTrainerSchedule(ctx context.Context, schedule trainer.TrainerSchedule) error {
	doc := TrainerScheduleDocument{
		UUID:          schedule.UUID(),
		TrainerUUID:   schedule.TrainerUUID(),
		Limit:         schedule.Limit(),
		CustomerUUIDs: schedule.CustomerUUIDs(),
		Name:          schedule.Name(),
		Desc:          schedule.Desc(),
		Date:          schedule.Date().Format(m.cfg.Format),
	}
	filter := bson.M{"_id": schedule.UUID(), "trainer_uuid": schedule.TrainerUUID()}
	err := m.updateOne(ctx, filter, doc)
	if err != nil {
		return fmt.Errorf("update one failed: %v", err)
	}
	return nil
}

func (m *MongoDB) QueryTrainerSchedule(ctx context.Context, UUID, trainerUUID string) (trainer.TrainerSchedule, error) {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
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
	date, err := time.Parse(m.cfg.Format, dst.Date)
	if err != nil {
		return trainer.TrainerSchedule{}, fmt.Errorf("parsing date value from document failed: %v", err)
	}
	out := trainer.UnmarshalFromDatabase(dst.UUID, dst.TrainerUUID, dst.Name, dst.Desc, dst.CustomerUUIDs, date, dst.Limit)
	return out, nil
}

func (m *MongoDB) QueryTrainerSchedules(ctx context.Context, trainerUUID string) ([]trainer.TrainerSchedule, error) {
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
		date, err := time.Parse(m.cfg.Format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		s := trainer.UnmarshalFromDatabase(d.UUID, d.TrainerUUID, d.Name, d.Desc, d.CustomerUUIDs, date, d.Limit)
		out = append(out, s)
	}
	return out, nil
}

func (m *MongoDB) CancelTrainerSchedules(ctx context.Context, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"trainer_uuid": trainerUUID}
	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return fmt.Errorf("delete many failed: %v", err)
	}
	return nil
}

func (m *MongoDB) CancelTrainerSchedule(ctx context.Context, UUID, trainerUUID string) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": UUID, "trainer_uuid": trainerUUID}
	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return fmt.Errorf("delete one failed: %v", err)
	}
	return nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.cli.Disconnect(ctx)
}

func (m *MongoDB) DropDatabase(ctx context.Context) error {
	db := m.cli.Database(m.cfg.Database)
	return db.Drop(ctx)
}

func (m *MongoDB) Ping(ctx context.Context) error {
	err := m.cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("ping request failed: %v", err)
	}
	return err
}

func (m *MongoDB) updateOne(ctx context.Context, filter bson.M, doc interface{}) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	update := bson.M{
		"$set": doc,
	}
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("update document failed: %v", err)
	}
	return nil
}

func (m *MongoDB) DropCollection(ctx context.Context) error {
	db := m.cli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	return coll.Drop(ctx)
}

func NewMongoDB(cfg MongoDBConfig) (*MongoDB, error) {
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

	m := MongoDB{
		cli: cli,
		cfg: cfg,
	}
	return &m, nil
}
