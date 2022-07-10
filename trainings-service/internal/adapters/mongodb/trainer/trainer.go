package trainer

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/common"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type RepositoryConfig struct {
	Addr              string
	Database          string
	TrainerCollection string
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
	Format            string
}

type Repository struct {
	cfg      RepositoryConfig
	commands *CommandHandler
	queries  *QueryHandler
}

func (r *Repository) UpsertTrainerWorkoutGroup(ctx context.Context, workoutGroup trainer.WorkoutGroup) error {
	return r.commands.UpsertTrainerWorkoutGroup(ctx, workoutGroup)
}
func (r *Repository) QueryTrainerWorkoutGroup(ctx context.Context, groupUUID, trainerUUID string) (trainer.WorkoutGroup, error) {
	return r.queries.QueryTrainerWorkoutGroup(ctx, groupUUID, trainerUUID)
}

func (r *Repository) QueryTrainerWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	return r.queries.QueryTrainerWorkoutGroups(ctx, trainerUUID)
}

func (r *Repository) QueryTrainerWorkoutGroupWithDate(ctx context.Context, trainerUUID string, date time.Time) (trainer.WorkoutGroup, error) {
	return r.queries.QueryTrainerWorkoutGroupWithDate(ctx, trainerUUID, date)
}

func (r *Repository) QueryCustomerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID, customerUUID string) (trainer.WorkoutGroup, error) {
	return r.queries.QueryCustomerWorkoutGroup(ctx, trainerUUID, groupUUID, customerUUID)
}

func (r *Repository) DeleteTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error {
	return r.commands.DeleteTrainerWorkoutGroups(ctx, trainerUUID)
}

func (r *Repository) DeleteTrainerWorkoutGroup(ctx context.Context, trainerUUID, groupUUID string) error {
	return r.commands.DeleteTrainerWorkoutGroup(ctx, trainerUUID, groupUUID)
}

func NewTrainerRepository(cfg RepositoryConfig) (*Repository, error) {
	mongoCLI, err := common.NewMongoClient(cfg.Addr, cfg.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("creating mongo cli failed: %v", err)
	}
	t := Repository{
		cfg: cfg,
		commands: NewCommandHandler(mongoCLI, CommandHandlerConfig{
			Collection:     cfg.TrainerCollection,
			Database:       cfg.Database,
			Format:         cfg.Format,
			CommandTimeout: cfg.CommandTimeout,
		}),
		queries: NewQueryHandler(mongoCLI, QueryHandlerConfig{
			Collection:   cfg.TrainerCollection,
			Database:     cfg.Database,
			Format:       cfg.Format,
			QueryTimeout: cfg.QueryTimeout,
		}),
	}
	return &t, nil
}
