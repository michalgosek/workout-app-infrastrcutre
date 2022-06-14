package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/testutil"
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

func (m *CustomerTestSuite) BeforeTest(mName, testName string) {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.ConnectionTimeout)
	defer cancel()
	mongoCLI, err := mongo.NewClient(options.Client().ApplyURI(m.cfg.Addr))
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
	m.mongoCLI = mongoCLI
	m.commandHandler = mongodb.NewCustomerCommandHandler(mongoCLI, mongodb.CustomerCommandHandlerConfig{
		Collection:     m.cfg.CustomerCollection,
		Database:       m.cfg.Database,
		Format:         m.cfg.Format,
		CommandTimeout: m.cfg.CommandTimeout,
	})
	err = m.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}

	m.queryHandler = mongodb.NewCustomerQueryHandler(mongoCLI, mongodb.CustomerQueryHandlerConfig{
		Collection:   m.cfg.CustomerCollection,
		Database:     m.cfg.Database,
		Format:       m.cfg.Format,
		QueryTimeout: m.cfg.QueryTimeout,
	})
}

func (m *CustomerTestSuite) AfterTest(suiteName, testName string) {
	ctx := context.Background()
	err := m.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}
}

func (m *CustomerTestSuite) TearDownSuite() {
	t := m.T()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, m.cfg.CommandTimeout)
	defer cancel()
	db := m.mongoCLI.Database(m.cfg.Database)
	err := db.Drop(ctx)
	if err != nil {
		t.Logf("mongoCLI cli dropping db failed: %v", err)
	}
	err = m.mongoCLI.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func (m *CustomerTestSuite) TestShouldInsertCustomerScheduleWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateCustomerSchedule(customerUUID)

	// when:
	err := m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QuerySchedule(ctx, expectedSchedule.CustomerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *CustomerTestSuite) TestShouldUpdateScheduelsListOfExistingCustomerSchedule() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateCustomerSchedule(trainerUUID)

	m.commandHandler.UpsertSchedule(ctx, expectedSchedule)
	expectedSchedule.AssignSchedule("371c59fb-588e-4335-8fb6-3b0186795fd0")

	// when:
	err := m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QuerySchedule(ctx, expectedSchedule.CustomerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *CustomerTestSuite) TestShouldDeleteCustomerScheduleWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const customerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const scheduleUUID = "2d0d3ce5-3ec4-48b1-a03b-a2fa5440963d"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateCustomerSchedule(customerUUID)
	expectedSchedule.AssignSchedule(scheduleUUID)

	m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// when:
	err := m.commandHandler.DeleteSchedule(ctx, customerUUID, expectedSchedule.UUID())

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QuerySchedule(ctx, expectedSchedule.CustomerUUID())
	assert.Nil(err)
	assert.Empty(actualSchedule)
}
