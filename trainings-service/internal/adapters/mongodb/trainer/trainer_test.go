package trainer_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

type Config struct {
	Addr              string
	Database          string
	TrainerCollection string
	CommandTimeout    time.Duration
	QueryTimeout      time.Duration
	ConnectionTimeout time.Duration
	Format            string
}

type TrainerTestSuite struct {
	suite.Suite
	commandHandler *trainer.CommandHandler
	queryHandler   *trainer.QueryHandler
	cfg            Config
	mongoCLI       *mongo.Client
}

func TestTrainerTestSuite_Integration(t *testing.T) {
	cfg := Config{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service",
		TrainerCollection: "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Format:            "2006-01-02 15:04",
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
	suite.Run(t, &TrainerTestSuite{
		mongoCLI: mongoCLI,
		commandHandler: trainer.NewCommandHandler(mongoCLI, trainer.CommandHandlerConfig{
			Collection:     cfg.TrainerCollection,
			Database:       cfg.Database,
			Format:         cfg.Format,
			CommandTimeout: cfg.CommandTimeout,
		}),
		queryHandler: trainer.NewQueryHandler(mongoCLI, trainer.QueryHandlerConfig{
			Collection:   cfg.TrainerCollection,
			Database:     cfg.Database,
			Format:       cfg.Format,
			QueryTimeout: cfg.QueryTimeout,
		}),
		cfg: cfg,
	})
}

func (m *TrainerTestSuite) AfterTest(string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.ConnectionTimeout)
	defer cancel()
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

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroup() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const groupUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const trainerUUID = "f7cb1930-11e4-490a-9e92-70be45aef11a"
	ctx := context.Background()

	// when:
	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
}

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroups() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const groupUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroups(ctx, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
}

func (m *TrainerTestSuite) TestShouldReturnWorkoutGroup() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// when:
	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
}

func (m *TrainerTestSuite) TestShouldInsertNonExistingWorkoutGroup() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)

	// when:
	err := m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
}

func (m *TrainerTestSuite) TestShouldUpdateNameOfExistingWorkoutGroup() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)
	_ = expectedGroup.UpdateDescription("dummy2")

	// when:
	err := m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
}

func (m *TrainerTestSuite) TestShouldDeleteWorkoutGroupWithSuccess() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// when:
	err := m.commandHandler.DeleteTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())
	assertions.Nil(err)
	assertions.Empty(actualGroup)
}

func (m *TrainerTestSuite) TestShouldNotReturnErrorWhenDeleteWorkoutGroupNonExist() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const fakeUUID = "13dd31ee-e131-44e1-8d95-dd6317af81b7"
	ctx := context.Background()
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// when:
	err := m.commandHandler.DeleteTrainerWorkoutGroup(ctx, fakeUUID, trainerUUID)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroup(ctx, trainerUUID, expectedGroup.UUID())
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
}

func (m *TrainerTestSuite) TestShouldDeleteAllTrainerWorkoutGroupsWithSuccess() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.NewTrainerWorkoutGroup(trainerUUID)
	second := testutil.NewTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, first)
	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, second)

	// when:
	err := m.commandHandler.DeleteTrainerWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualGroups, err := m.queryHandler.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	assertions.Nil(err)
	assertions.Empty(actualGroups)
}

func (m *TrainerTestSuite) TestShouldNotReturnErrorWhenDeleteWorkoutGroupsNonExist() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	err := m.commandHandler.DeleteTrainerWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualGroups, err := m.queryHandler.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	assertions.Nil(err)
	assertions.Empty(actualGroups)
}

func (m *TrainerTestSuite) TestShouldReturnWorkoutGroupWithSpecifiedDate() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)
	_ = m.commandHandler.UpsertTrainerWorkoutGroup(ctx, expectedGroup)

	// when:
	actualGroup, err := m.queryHandler.QueryTrainerWorkoutGroupWithDate(ctx, trainerUUID, expectedGroup.Date())

	// then:
	assertions.Nil(err)
	assertions.Equal(actualGroup, expectedGroup)
}
