package mongodb

import (
	"context"
	"fmt"
	"time"

	mcustomer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
)

type customerHandler struct {
	commands *mcustomer.CommandHandler
	queries  *mcustomer.QueryHandler
}

func (c *customerHandler) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	return c.commands.UpsertCustomerWorkoutDay(ctx, workout)
}

func (c *customerHandler) QueryCustomerWorkoutDay(ctx context.Context, customerUUID, GroupUUID string) (customer.WorkoutDay, error) {
	return c.queries.QueryCustomerWorkoutDay(ctx, customerUUID, GroupUUID)
}

func (c *customerHandler) QueryCustomerWorkoutDays(ctx context.Context, customerUUID string) ([]customer.WorkoutDay, error) {
	return c.queries.QueryCustomerWorkoutDays(ctx, customerUUID)
}

func (c *customerHandler) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID, workoutDayUUID string) error {
	return c.commands.DeleteCustomerWorkoutDay(ctx, customerUUID, workoutDayUUID)
}

func (c *customerHandler) DeleteCustomerWorkoutDays(ctx context.Context, customerUUID string) error {
	return c.commands.DeleteCustomerWorkoutDays(ctx, customerUUID)
}

type CustomerRepositoryConfig struct {
	Addr               string
	Database           string
	CustomerCollection string
	CommandTimeout     time.Duration
	QueryTimeout       time.Duration
	ConnectionTimeout  time.Duration
	Format             string
}

type CustomerRepository struct {
	cfg CustomerRepositoryConfig
	customerHandler
}

func NewCustomerRepository(cfg CustomerRepositoryConfig) (*CustomerRepository, error) {
	mongoCLI, err := newMongoClient(cfg.Addr, cfg.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("creating mongo cli failed: %v", err)
	}
	c := CustomerRepository{
		cfg: cfg,
		customerHandler: customerHandler{
			commands: mcustomer.NewCommandHandler(mongoCLI, mcustomer.CommandHandlerConfig{
				Collection:     cfg.CustomerCollection,
				Database:       cfg.Database,
				Format:         cfg.Format,
				CommandTimeout: cfg.CommandTimeout,
			}),
			queries: mcustomer.NewQueryHandler(mongoCLI, mcustomer.QueryHandlerConfig{
				Collection:   cfg.CustomerCollection,
				Database:     cfg.Database,
				Format:       cfg.Format,
				QueryTimeout: cfg.QueryTimeout,
			}),
		},
	}
	return &c, nil
}
