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

func (m *TrainerTestSuite) BeforeTest(string, string) {
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

func (m *TrainerTestSuite) AfterTest(string, string) {
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

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroupWhenQuery() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const groupUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
}

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroupsWhenQuery() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const groupUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualGroup, err := m.queryHandler.QueryWorkoutGroups(ctx, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualGroup)
}

func (m *TrainerTestSuite) TestShouldReturnWorkoutGroupForQuery() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	expectedGroup := testutil.NewTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)

	// when:
	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, expectedGroup.UUID())

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
	err := m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, expectedGroup.UUID())
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

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)
	_ = expectedGroup.UpdateDesc("dummy2")

	// when:
	err := m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, expectedGroup.UUID())
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
	_ = m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)

	// when:
	err := m.commandHandler.DeleteWorkoutGroup(ctx, expectedGroup.UUID())

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, expectedGroup.UUID())
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

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, expectedGroup)

	// when:
	err := m.commandHandler.DeleteWorkoutGroup(ctx, fakeUUID)

	// then:
	assertions.Nil(err)

	actualGroup, err := m.queryHandler.QueryWorkoutGroup(ctx, expectedGroup.UUID())
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

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, first)
	_ = m.commandHandler.UpsertWorkoutGroup(ctx, second)

	// when:
	err := m.commandHandler.DeleteWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualGroups, err := m.queryHandler.QueryWorkoutGroups(ctx, trainerUUID)
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
	err := m.commandHandler.DeleteWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualGroups, err := m.queryHandler.QueryWorkoutGroups(ctx, trainerUUID)
	assertions.Nil(err)
	assertions.Empty(actualGroups)
}
