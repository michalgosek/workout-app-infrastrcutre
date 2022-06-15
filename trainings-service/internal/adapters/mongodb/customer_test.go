package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type CustomerTestSuite struct {
	suite.Suite
	commandHandler *mongodb.CustomerCommandHandler
	queryHandler   *mongodb.CustomerQueryHandler
	cfg            mongodb.Config
	mongoCLI       *mongo.Client
}

func TestCustomerTestSuite_Integration(t *testing.T) {
	suite.Run(t, &CustomerTestSuite{
		cfg: mongodb.Config{
			Addr:               "mongodb://localhost:27017",
			Database:           "trainings_service_test",
			CustomerCollection: "customer_schedules",
			CommandTimeout:     10 * time.Second,
			QueryTimeout:       10 * time.Second,
			ConnectionTimeout:  10 * time.Second,
			Format:             "2006-01-02 15:04",
		},
	})
}

func (c *CustomerTestSuite) BeforeTest(string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.ConnectionTimeout)
	defer cancel()
	mongoCLI, err := mongo.NewClient(options.Client().ApplyURI(c.cfg.Addr))
	if err != nil {
		panic(err)
	}
	err = mongoCLI.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = mongoCLI.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	c.mongoCLI = mongoCLI
	c.commandHandler = mongodb.NewCustomerCommandHandler(mongoCLI, mongodb.CustomerCommandHandlerConfig{
		Collection:     c.cfg.CustomerCollection,
		Database:       c.cfg.Database,
		Format:         c.cfg.Format,
		CommandTimeout: c.cfg.CommandTimeout,
	})
	err = c.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}

	c.queryHandler = mongodb.NewCustomerQueryHandler(mongoCLI, mongodb.CustomerQueryHandlerConfig{
		Collection:   c.cfg.CustomerCollection,
		Database:     c.cfg.Database,
		Format:       c.cfg.Format,
		QueryTimeout: c.cfg.QueryTimeout,
	})
}

func (c *CustomerTestSuite) AfterTest(string, string) {
	ctx := context.Background()
	err := c.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}
}

func (c *CustomerTestSuite) TearDownSuite() {
	t := c.T()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.cfg.CommandTimeout)
	defer cancel()
	db := c.mongoCLI.Database(c.cfg.Database)
	err := db.Drop(ctx)
	if err != nil {
		t.Logf("mongoCLI cli dropping db failed: %v", err)
	}
	err = c.mongoCLI.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func (c *CustomerTestSuite) TestShouldUpsertCustomerWorkoutDayWhenNotExistWithSuccess() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workout := testutil.NewWorkoutDay(customerUUID)

	// when:
	err := c.commandHandler.UpsertCustomerWorkoutDay(ctx, workout)

	// then:
	assertions.Nil(err)

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.TrainerWorkoutGroupUUID())
	assertions.Nil(err)
	assertions.Equal(workout, actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldDeleteCustomerWorkoutDayWithSuccess() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workout := testutil.NewWorkoutDay(customerUUID)

	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, workout)

	// when:
	err := c.commandHandler.DeleteCustomerWorkoutDay(ctx, customerUUID, workout.UUID())

	// then:
	assertions.Nil(err)

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.TrainerWorkoutGroupUUID())
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDay)
}
