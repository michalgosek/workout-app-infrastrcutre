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

type TrainerTestSuite struct {
	suite.Suite
	commandHandler *mongodb.TrainerCommandHandler
	queryHandler   *mongodb.TrainerQueryHandler
	cfg            mongodb.Config
	mongoCLI       *mongo.Client
}

func TestTrainerTestSuite_Integration(t *testing.T) {
	suite.Run(t, &TrainerTestSuite{
		cfg: mongodb.Config{
			Addr:              "mongodb://localhost:27017",
			Database:          "trainings_service_test",
			TrainerCollection: "trainer_schedules",
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
			Format:            "2006-01-02 15:04",
		},
	})
}

func (m *TrainerTestSuite) BeforeTest(mName, testName string) {
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
	m.commandHandler = mongodb.NewTrainerCommandHandler(mongoCLI, mongodb.TrainerCommandHandlerConfig{
		Collection:     m.cfg.TrainerCollection,
		Database:       m.cfg.Database,
		Format:         m.cfg.Format,
		CommandTimeout: m.cfg.CommandTimeout,
	})
	err = m.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}

	m.queryHandler = mongodb.NewTrainerQueryHandler(mongoCLI, mongodb.TrainerQueryHandlerConfig{
		Collection:   m.cfg.TrainerCollection,
		Database:     m.cfg.Database,
		Format:       m.cfg.Format,
		QueryTimeout: m.cfg.CommandTimeout,
	})
}

func (m *TrainerTestSuite) AfterTest(suiteName, testName string) {
	ctx := context.Background()
	err := m.commandHandler.DropCollection(ctx)
	if err != nil {
		panic(err)
	}
}

func (m *TrainerTestSuite) TearDownSuite() {
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

func (m *TrainerTestSuite) TestShouldInsertTrainerScheduleWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	// when:
	err := m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldUpdateNameOfExisitingSchedule() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.commandHandler.UpsertSchedule(ctx, expectedSchedule)
	expectedSchedule.UpdateDesc("dummy2")

	// when:
	err := m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldCancelTrainerScheduleWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// when:
	err := m.commandHandler.CancelSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Empty(actualSchedule)
}

func (m *TrainerTestSuite) TestShouldReturnEmptyScheduleWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const fakeUUID = "13dd31ee-e131-44e1-8d95-dd6317af81b7"
	const fakeTrainerUUID = "b71f6ed1-a982-47f2-a3fe-1d32cf3a132f"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.commandHandler.UpsertSchedule(ctx, expectedSchedule)

	// when:
	err := m.commandHandler.CancelSchedule(ctx, fakeUUID, fakeTrainerUUID)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.queryHandler.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldCancelAllUpsertedSchedules() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)

	m.commandHandler.UpsertSchedule(ctx, first)
	m.commandHandler.UpsertSchedule(ctx, second)

	// when:
	err := m.commandHandler.CancelSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.queryHandler.QueryTrainerSchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

func (m *TrainerTestSuite) TestShouldReturnEmptySchedulesWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	err := m.commandHandler.CancelSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.queryHandler.QueryTrainerSchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}
