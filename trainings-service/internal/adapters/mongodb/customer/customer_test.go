package customer_test

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"

	mcustomer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Addr               string
	Database           string
	TrainerCollection  string
	CustomerCollection string
	CommandTimeout     time.Duration
	QueryTimeout       time.Duration
	ConnectionTimeout  time.Duration
	Format             string
}

type CustomerTestSuite struct {
	suite.Suite
	mongoCLI       *mongo.Client
	commandHandler *mcustomer.CommandHandler
	queryHandler   *mcustomer.QueryHandler
	cfg            Config
}

func TestCustomerTestSuite_Integration(t *testing.T) {
	cfg := Config{
		Addr:               "mongodb://localhost:27017",
		Database:           "trainings_service_test",
		CustomerCollection: "customer_schedules",
		CommandTimeout:     10 * time.Second,
		QueryTimeout:       10 * time.Second,
		ConnectionTimeout:  10 * time.Second,
		Format:             "2006-01-02 15:04",
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()
	mongoCLI, err := mongo.NewClient(options.Client().ApplyURI(cfg.Addr))
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

	suite.Run(t, &CustomerTestSuite{
		mongoCLI: mongoCLI,
		commandHandler: mcustomer.NewCommandHandler(mongoCLI, mcustomer.CommandHandlerConfig{
			Collection:     cfg.CustomerCollection,
			Database:       cfg.Database,
			Format:         cfg.Format,
			CommandTimeout: cfg.CommandTimeout,
		}),
		queryHandler: mcustomer.NewQueryHandler(mongoCLI, mcustomer.QueryHandlerConfig{
			Collection:   cfg.CustomerCollection,
			Database:     cfg.Database,
			Format:       cfg.Format,
			QueryTimeout: cfg.QueryTimeout,
		}),
		cfg: Config{
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

func (c *CustomerTestSuite) TestShouldReturnEmptyWorkoutDayForQuery() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const trainerUUID = "e1981243-0ac7-4dee-99c6-a8b7a0b13f95"
	ctx := context.Background()

	// when:
	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, customerUUID, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldReturnWorkoutDayForQuery() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	expectedWorkoutDay := testutil.NewWorkoutDay(customerUUID)
	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, expectedWorkoutDay)

	// when:
	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, customerUUID, expectedWorkoutDay.GroupUUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkoutDay, actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldReturnEmptyWorkoutDaysForQuery() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualWorkoutDays, err := c.queryHandler.QueryCustomerWorkoutDays(ctx, customerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)
}

func (c *CustomerTestSuite) TestShouldReturnWorkoutDaysForQuery() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	first := testutil.NewWorkoutDay(customerUUID)
	second := testutil.NewWorkoutDay(customerUUID)
	expectedWorkouts := []customer.WorkoutDay{first, second}

	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, first)
	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, second)

	// when:
	days, err := c.queryHandler.QueryCustomerWorkoutDays(ctx, customerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkouts, days)
}

func (c *CustomerTestSuite) TestShouldUpsertNewWorkoutDayWithSuccess() {
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

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.GroupUUID())
	assertions.Nil(err)
	assertions.Equal(workout, actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldNotReturnErrorWhenDeleteWorkoutDayNonExist() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workout := testutil.NewWorkoutDay(customerUUID)

	// when:
	err := c.commandHandler.DeleteCustomerWorkoutDay(ctx, customerUUID, workout.UUID())

	// then:
	assertions.Nil(err)

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.GroupUUID())
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldDeleteWorkoutDayWithSuccess() {
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

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.GroupUUID())
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDay)
}

func (c *CustomerTestSuite) TestShouldNotReturnErrorWhenDeleteWorkoutDaysNonExist() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workout := testutil.NewWorkoutDay(customerUUID)

	// when:
	err := c.commandHandler.DeleteCustomerWorkoutDays(ctx, customerUUID)

	// then:
	assertions.Nil(err)

	actualWorkoutDays, err := c.queryHandler.QueryCustomerWorkoutDays(ctx, workout.CustomerUUID())
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)
}

func (c *CustomerTestSuite) TestShouldDeleteWorkoutDaysWithSuccess() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workout := testutil.NewWorkoutDay(customerUUID)

	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, workout)

	// when:
	err := c.commandHandler.DeleteCustomerWorkoutDays(ctx, customerUUID)

	// then:
	assertions.Nil(err)

	actualWorkoutDays, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, workout.CustomerUUID(), workout.GroupUUID())
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)
}

func (c *CustomerTestSuite) TestShouldDeleteAllCustomerWorkoutDaysWithSpecifiedGroupUUID() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const (
		firstCustomerName  = "John Doe"
		secondCustomerName = "Jerry Doe"
		commonTrainerUUID  = "1fdae430-f4bf-47d9-bb4c-a161a7d5a304"
		firstCustomerUUID  = "f1741a08-39d7-465d-adc9-a63cf058b409"
		secondCustomerUUID = "2e9be2da-448d-4236-afb6-8bffba6668b6"
		groupUUID          = "0f2407e1-346c-4f56-bceb-d94c9d3b014b"
	)
	ctx := context.Background()
	schedule := time.Now().AddDate(0, 0, 1)
	firstCustomerWorkoutDay, _ := customer.NewWorkoutDay(firstCustomerUUID, firstCustomerName, groupUUID, commonTrainerUUID, schedule)
	secondCustomerWorkoutDay, _ := customer.NewWorkoutDay(secondCustomerUUID, secondCustomerName, groupUUID, commonTrainerUUID, schedule)

	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, firstCustomerWorkoutDay)
	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, secondCustomerWorkoutDay)

	// when:
	err := c.commandHandler.DeleteCustomersWorkoutDaysWithGroup(ctx, groupUUID)

	// then:
	assertions.Nil(err)

	actualWorkoutDays, err := c.queryHandler.QueryCustomerWorkoutDays(ctx, firstCustomerUUID)
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)

	actualWorkoutDays, err = c.queryHandler.QueryCustomerWorkoutDays(ctx, secondCustomerUUID)
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)
}

func (c *CustomerTestSuite) TestShouldDeleteTwoCustomerWorkoutDaysWithSpecifiedGroupUUID() {
	t := c.T()
	assertions := assert.New(t)

	// given:
	const (
		firstCustomerName  = "John Doe"
		secondCustomerName = "Jerry Doe"
		thirdCustomerName  = "Jane Doe"
		firstCustomerUUID  = "f1741a08-39d7-465d-adc9-a63cf058b409"
		secondCustomerUUID = "2e9be2da-448d-4236-afb6-8bffba6668b6"
		thirdCustomerUUID  = "1d25c898-dadc-4b98-a065-1dd7e0d4c904"
		commonGroupUUID    = "0f2407e1-346c-4f56-bceb-d94c9d3b014b"
		commonTrainerUUID  = "43404fd3-62d2-409e-8aa6-2f85ac290cdc"
		otherGroupUUID     = "ece96132-a4fd-466d-9c5c-965a4744f0ae"
	)
	ctx := context.Background()
	schedule, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	firstCustomerWorkoutDay, _ := customer.NewWorkoutDay(firstCustomerUUID, firstCustomerName, commonGroupUUID, commonTrainerUUID, schedule)
	secondCustomerWorkoutDay, _ := customer.NewWorkoutDay(secondCustomerUUID, secondCustomerName, commonGroupUUID, commonTrainerUUID, schedule)
	thirdCustomerWorkoutDay, _ := customer.NewWorkoutDay(thirdCustomerUUID, thirdCustomerName, otherGroupUUID, commonTrainerUUID, schedule)

	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, firstCustomerWorkoutDay)
	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, secondCustomerWorkoutDay)
	_ = c.commandHandler.UpsertCustomerWorkoutDay(ctx, thirdCustomerWorkoutDay)

	// when:
	err := c.commandHandler.DeleteCustomersWorkoutDaysWithGroup(ctx, commonGroupUUID)

	// then:
	assertions.Nil(err)

	actualWorkoutDays, err := c.queryHandler.QueryCustomerWorkoutDays(ctx, firstCustomerUUID)
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)

	actualWorkoutDays, err = c.queryHandler.QueryCustomerWorkoutDays(ctx, secondCustomerUUID)
	assertions.Nil(err)
	assertions.Empty(actualWorkoutDays)

	actualWorkoutDay, err := c.queryHandler.QueryCustomerWorkoutDay(ctx, thirdCustomerUUID, otherGroupUUID)
	assertions.Nil(err)
	assertions.Equal(actualWorkoutDay, thirdCustomerWorkoutDay)
}
