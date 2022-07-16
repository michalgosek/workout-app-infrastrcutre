package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Timeouts struct {
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
}

type Config struct {
	Addr       string
	Database   string
	Collection string
	Timeouts   Timeouts
}

type Repository struct {
	cli *mongo.Client
	cfg Config
}

func (r *Repository) InsertTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	doc := WorkoutGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: TrainerWorkoutGroupsWriteModel{
			UUID: g.Trainer().UUID(),
			Name: g.Trainer().Name(),
		},
		Limit: g.Limit(),
	}
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateTrainerWorkoutGroup(ctx context.Context, g *trainings.WorkoutGroup) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	doc := WorkoutGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: TrainerWorkoutGroupsWriteModel{
			UUID: g.Trainer().UUID(),
			Name: g.Trainer().Name(),
		},
		Limit:        g.Limit(),
		Participants: convertToWriteModelParticipants(g.Participants()...),
	}
	filter := bson.M{"_id": g.UUID(), "trainer._id": g.Trainer().UUID()}
	update := bson.M{"$set": doc}
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) error {
	f := bson.M{"_id": groupUUID, "trainer._id": trainerUUID}
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	_, err := coll.DeleteOne(ctx, f)
	if err != nil {
		return nil
	}
	return nil
}

func (r *Repository) TrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (query.TrainerWorkoutGroup, error) {
	f := bson.M{"_id": groupUUID, "trainer._id": trainerUUID}
	g, err := r.findTrainerWorkoutGroupWithFilter(ctx, f)
	if err != nil {
		return query.TrainerWorkoutGroup{}, nil
	}
	m := UnmarshalToQueryTrainerWorkoutGroup(g)
	return m, nil
}

func (r *Repository) TrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]query.TrainerWorkoutGroup, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	f := bson.M{"trainer._id": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dd []WorkoutGroupWriteModel
	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}
	m := UnmarshalToQueryTrainerWorkoutGroups(dd...)
	return m, nil
}

func (r *Repository) IsDuplicateTrainerWorkoutGroupExists(ctx context.Context, groupUUID, trainerUUID string) (bool, error) {
	f := bson.M{"_id": groupUUID, "trainer._id": trainerUUID}
	doc, err := r.findTrainerWorkoutGroupWithFilter(ctx, f)
	if err != nil {
		return false, err
	}
	return doc.UUID != "", nil
}

func (r *Repository) QueryTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainings.WorkoutGroup, error) {
	f := bson.M{"_id": groupUUID, "trainer._id": trainerUUID}
	doc, err := r.findTrainerWorkoutGroupWithFilter(ctx, f)
	if err != nil {
		return trainings.WorkoutGroup{}, err
	}

	var pp []trainings.DatabaseWorkoutGroupParticipant
	for _, p := range doc.Participants {
		pp = append(pp, trainings.DatabaseWorkoutGroupParticipant{UUID: p.UUID, Name: p.Name})
	}
	g := trainings.UnmarshalWorkoutGroupFromDatabase(trainings.DatabaseWorkoutGroup{
		UUID:        doc.UUID,
		Name:        doc.Name,
		Description: doc.Description,
		Limit:       doc.Limit,
		Date:        doc.Date,
		Trainer: trainings.DatabaseWorkoutGroupTrainer{
			UUID: doc.Trainer.UUID,
			Name: doc.Trainer.Name,
		},
		Participants: pp,
	})
	return g, nil
}

func (r *Repository) findTrainerWorkoutGroupWithFilter(ctx context.Context, f bson.M) (WorkoutGroupWriteModel, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
	defer cancel()
	res := coll.FindOne(ctx, f)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return WorkoutGroupWriteModel{}, nil
	}
	if res.Err() != nil {
		return WorkoutGroupWriteModel{}, fmt.Errorf("find one failed: %s", res.Err())
	}

	var doc WorkoutGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return WorkoutGroupWriteModel{}, fmt.Errorf("decode failed: %s", res.Err())
	}
	return doc, nil
}

func NewRepository(cfg Config) (*Repository, error) {
	cli, err := NewClient(cfg.Addr, cfg.Timeouts.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("mongo cli creation failed: %w", err)
	}
	m := Repository{
		cli: cli,
		cfg: cfg,
	}
	return &m, nil
}

func NewClient(addr string, d time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(d)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("mongo client creation failed: %v", err)
	}
	err = cli.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("mongo client connection failed: %v", err)
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("mongo client ping req failed: %v", err)
	}
	return cli, nil
}

func convertToWriteModelParticipants(pp ...trainings.Participant) []ParticipantWriteModel {
	var out []ParticipantWriteModel
	for _, p := range pp {
		out = append(out, ParticipantWriteModel{
			UUID: p.UUID(),
			Name: p.Name(),
		})
	}
	return out
}
