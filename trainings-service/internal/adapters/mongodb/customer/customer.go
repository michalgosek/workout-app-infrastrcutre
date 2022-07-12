package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/common"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
)

type RepositoryConfig struct {
	Addr               string
	Database           string
	CustomerCollection string
	CommandTimeout     time.Duration
	QueryTimeout       time.Duration
	ConnectionTimeout  time.Duration
	Format             string
}

type Repository struct {
	cfg      RepositoryConfig
	commands *CommandHandler
	queries  *QueryHandler
}

func (r *Repository) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	return r.commands.UpsertCustomerWorkoutDay(ctx, workout)
}

func (r *Repository) QueryCustomerWorkoutDay(ctx context.Context, customerUUID, GroupUUID string) (customer.WorkoutDay, error) {
	return r.queries.QueryCustomerWorkoutDay(ctx, customerUUID, GroupUUID)
}

func (r *Repository) QueryCustomerWorkoutDays(ctx context.Context, customerUUID string) ([]customer.WorkoutDay, error) {
	return r.queries.QueryCustomerWorkoutDays(ctx, customerUUID)
}

func (r *Repository) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, workoutDayUUID string) error {
	return r.commands.DeleteCustomerWorkoutDay(ctx, customerUUID, workoutDayUUID)
}

func (r *Repository) DeleteCustomersWorkoutDaysWithGroup(ctx context.Context, groupUUID string) error {
	return r.commands.DeleteCustomersWorkoutDaysWithGroup(ctx, groupUUID)
}

func (r *Repository) DeleteCustomerWorkoutDays(ctx context.Context, customerUUID string) error {
	return r.commands.DeleteCustomerWorkoutDays(ctx, customerUUID)
}

func (r *Repository) DeleteCustomersWorkoutDaysWithTrainer(ctx context.Context, trainerUUID string) error {
	return r.commands.DeleteCustomersWorkoutDaysWithTrainer(ctx, trainerUUID)
}

func NewCustomerRepository(cfg RepositoryConfig) (*Repository, error) {
	mongoCLI, err := common.NewMongoClient(cfg.Addr, cfg.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("creating mongo cli failed: %v", err)
	}
	c := Repository{
		cfg: cfg,
		commands: NewCommandHandler(mongoCLI, CommandHandlerConfig{
			Collection:     cfg.CustomerCollection,
			Database:       cfg.Database,
			Format:         cfg.Format,
			CommandTimeout: cfg.CommandTimeout,
		}),
		queries: NewQueryHandler(mongoCLI, QueryHandlerConfig{
			Collection:   cfg.CustomerCollection,
			Database:     cfg.Database,
			Format:       cfg.Format,
			QueryTimeout: cfg.QueryTimeout,
		}),
	}
	return &c, nil
}
