package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *Repository) findTrainingGroupWithFilter(ctx context.Context, f bson.M) (TrainingGroupWriteModel, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return TrainingGroupWriteModel{}, res.Err()
	}

	var doc TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return TrainingGroupWriteModel{}, err
	}
	return doc, nil
}

func (r *Repository) InsertTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	doc := TrainingGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: TrainerWriteModel{
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

func (r *Repository) UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.QueryTimeout)
	defer cancel()
	_, err := r.QueryTrainingGroup(ctx, g.UUID())
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	}
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	doc := TrainingGroupWriteModel{
		UUID:        g.UUID(),
		Name:        g.Name(),
		Description: g.Description(),
		Date:        g.Date(),
		Trainer: TrainerWriteModel{
			UUID: g.Trainer().UUID(),
			Name: g.Trainer().Name(),
		},
		Limit:        g.Limit(),
		Participants: ConvertToWriteModelParticipants(g.Participants()...),
	}
	filter := bson.M{"_id": g.UUID(), "trainer._id": g.Trainer().UUID()}
	update := bson.M{"$set": doc}
	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsTrainingGroupDuplicated(ctx context.Context, g *trainings.TrainingGroup) (bool, error) {
	f := bson.M{"trainer._id": g.Trainer().UUID(), "description": g.Description(), "name": g.Name(), "date": g.Date()}
	_, err := r.findTrainingGroupWithFilter(ctx, f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) DeleteTrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) error {
	f := bson.M{"_id": trainingUUID, "trainer._id": trainerUUID}
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

func (r *Repository) DeleteTrainingGroups(ctx context.Context, trainerUUID string) error {
	f := bson.M{"trainer._id": trainerUUID}
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)

	_, err := coll.DeleteMany(ctx, f)
	if err != nil {
		return nil
	}
	return nil
}

func (r *Repository) QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error) {
	f := bson.M{"_id": trainingUUID}
	doc, err := r.findTrainingGroupWithFilter(ctx, f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return trainings.TrainingGroup{}, nil
	}
	if err != nil {
		return trainings.TrainingGroup{}, err
	}

	var pp []trainings.DatabaseTrainingGroupParticipant
	for _, p := range doc.Participants {
		pp = append(pp, trainings.DatabaseTrainingGroupParticipant{UUID: p.UUID, Name: p.Name})
	}
	g := trainings.UnmarshalTrainingGroupFromDatabase(trainings.DatabaseTrainingGroup{
		UUID:        doc.UUID,
		Name:        doc.Name,
		Description: doc.Description,
		Limit:       doc.Limit,
		Date:        doc.Date,
		Trainer: trainings.DatabaseTrainingGroupTrainer{
			UUID: doc.Trainer.UUID,
			Name: doc.Trainer.Name,
		},
		Participants: pp,
	})
	return g, nil
}

func (r *Repository) TrainingGroup(ctx context.Context, trainingUUID, trainerUUID string) (query.TrainerWorkoutGroup, error) {
	f := bson.M{"_id": trainingUUID, "trainer._id": trainerUUID}
	g, err := r.findTrainingGroupWithFilter(ctx, f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return query.TrainerWorkoutGroup{}, nil
	}
	if err != nil {
		return query.TrainerWorkoutGroup{}, err
	}
	m := UnmarshalToQueryTrainerWorkoutGroup(g)
	return m, nil
}

func (r *Repository) TrainingGroups(ctx context.Context, trainerUUID string) ([]query.TrainerWorkoutGroup, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	ctx, cancel := context.WithTimeout(ctx, r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	f := bson.M{"trainer._id": trainerUUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dd []TrainingGroupWriteModel
	err = cur.All(ctx, &dd)
	if err != nil {
		return nil, err
	}
	m := UnmarshalToQueryTrainerWorkoutGroups(dd...)
	return m, nil
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
