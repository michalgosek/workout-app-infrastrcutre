package mongodb

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	readm "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Commands struct {
	*command.InsertTrainerGroupHandler
	*command.UpdateTrainingGroupHandler
	*command.DeleteTrainingGroupHandler
	*command.DeleteTrainingGroupsHandler
}

type Queries struct {
	*query.TrainerGroupHandler
	*query.TrainerGroupsHandler
	*query.AllTrainingsHandler
	*query.ParticipantGroupsHandler
}

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
	cfg      Config
	cli      *mongo.Client
	commands *Commands
	queries  *Queries
}

func (r *Repository) ParticipantGroups(ctx context.Context, UUID string) ([]readm.ParticipantGroup, error) {
	return r.queries.ParticipantGroupsHandler.Do(ctx, UUID)
}

func (r *Repository) findTrainingGroupWithFilter(ctx context.Context, f bson.M) (documents.TrainingGroupWriteModel, error) {
	db := r.cli.Database(r.cfg.Database)
	coll := db.Collection(r.cfg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.QueryTimeout)
	defer cancel()

	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return documents.TrainingGroupWriteModel{}, res.Err()
	}

	var doc documents.TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return documents.TrainingGroupWriteModel{}, err
	}
	return doc, nil
}

func (r *Repository) AllTrainingGroups(ctx context.Context) ([]readm.TrainingGroup, error) {
	return r.queries.AllTrainingsHandler.Do(ctx)
}

func (r *Repository) TrainerGroups(ctx context.Context, trainerUUID string) ([]readm.TrainerGroup, error) {
	return r.queries.TrainerGroupsHandler.Do(ctx, trainerUUID)
}

func (r *Repository) TrainerGroup(ctx context.Context, trainingUUID, trainerUUID string) (readm.TrainerGroup, error) {
	return r.queries.TrainerGroupHandler.Do(ctx, trainingUUID, trainerUUID)
}

func (r *Repository) InsertTrainerGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	return r.commands.InsertTrainerGroupHandler.Do(ctx, g)
}

func (r *Repository) UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	return r.commands.UpdateTrainingGroupHandler.Do(ctx, g)
}

func (r *Repository) DeleteTrainerGroup(ctx context.Context, trainingUUID, trainerUUID string) error {
	return r.commands.DeleteTrainingGroupHandler.Do(ctx, trainingUUID, trainerUUID)
}

func (r *Repository) DeleteTrainerGroups(ctx context.Context, trainerUUID string) error {
	return r.commands.DeleteTrainingGroupsHandler.Do(ctx, trainerUUID)
}

func (r Repository) QueryTrainingGroup(ctx context.Context, trainingUUID string) (trainings.TrainingGroup, error) {
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

func (r Repository) IsTrainingGroupDuplicated(ctx context.Context, g *trainings.TrainingGroup) (bool, error) {
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

func (r *Repository) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeouts.CommandTimeout)
	defer cancel()
	db := r.cli.Database(r.cfg.Database)
	err := db.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(cfg Config) (*Repository, error) {
	cli, err := NewClient(cfg.Addr, cfg.Timeouts.ConnectionTimeout)
	if err != nil {
		return nil, err
	}
	commandCfg := command.Config{
		Database:       cfg.Database,
		Collection:     cfg.Collection,
		CommandTimeout: cfg.Timeouts.CommandTimeout,
	}
	queryCfg := query.Config{
		Database:     cfg.Database,
		Collection:   cfg.Collection,
		QueryTimeout: cfg.Timeouts.QueryTimeout,
	}
	r := Repository{
		cfg: cfg,
		cli: cli,
		commands: &Commands{
			InsertTrainerGroupHandler:   command.NewInsertTrainerGroupHandler(cli, commandCfg),
			DeleteTrainingGroupHandler:  command.NewDeleteTrainingGroupHandler(cli, commandCfg),
			DeleteTrainingGroupsHandler: command.NewDeleteTrainingGroupsHandler(cli, commandCfg),
			UpdateTrainingGroupHandler:  command.NewUpdateTrainingGroupHandler(cli, commandCfg),
		},
		queries: &Queries{
			TrainerGroupHandler:      query.NewTrainerGroupHandler(cli, queryCfg),
			TrainerGroupsHandler:     query.NewTrainerGroupsHandler(cli, queryCfg),
			AllTrainingsHandler:      query.NewAllTrainingsHandler(cli, queryCfg),
			ParticipantGroupsHandler: query.NewParticipantGroupsHandler(cli, queryCfg),
		},
	}
	return &r, nil
}

func NewClient(addr string, d time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(d)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	err = cli.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return cli, nil
}
