package mongodb

import (
	"context"
	"fmt"
	"time"

	mtrainer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type trainerHandler struct {
	commands *mtrainer.CommandHandler
	queries  *mtrainer.QueryHandler
}

func (m *trainerHandler) UpsertTrainerWorkoutGroup(ctx context.Context, workoutGroup trainer.WorkoutGroup) error {
	return m.commands.UpsertTrainerWorkoutGroup(ctx, workoutGroup)
}
func (m *trainerHandler) QueryTrainerWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error) {
	return m.queries.QueryTrainerWorkoutGroup(ctx, groupUUID)
}

func (m *trainerHandler) QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	return m.queries.QueryTrainerWorkoutGroups(ctx, trainerUUID)
}

func (m *trainerHandler) DeleteTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error {
	return m.commands.DeleteTrainerWorkoutGroups(ctx, trainerUUID)
}

func (m *trainerHandler) DeleteTrainerWorkoutGroup(ctx context.Context, groupUUID string) error {
	return m.commands.DeleteTrainerWorkoutGroup(ctx, groupUUID)
}

type TrainerRepositoryConfig struct {
	Addr              string
	Database          string
	TrainerCollection string
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
	Format            string
}

type TrainerRepository struct {
	cfg TrainerRepositoryConfig
	trainerHandler
}

func NewTrainerRepository(cfg TrainerRepositoryConfig) (*TrainerRepository, error) {
	mongoCLI, err := newMongoClient(cfg.Addr, cfg.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("creating mongo cli failed: %v", err)
	}
	t := TrainerRepository{
		cfg: cfg,
		trainerHandler: trainerHandler{
			commands: mtrainer.NewCommandHandler(mongoCLI, mtrainer.CommandHandlerConfig{
				Collection:     cfg.TrainerCollection,
				Database:       cfg.Database,
				Format:         cfg.Format,
				CommandTimeout: cfg.CommandTimeout,
			}),
			queries: mtrainer.NewQueryHandler(mongoCLI, mtrainer.QueryHandlerConfig{
				Collection:   cfg.TrainerCollection,
				Database:     cfg.Database,
				Format:       cfg.Format,
				QueryTimeout: cfg.QueryTimeout,
			}),
		},
	}
	return &t, nil
}
