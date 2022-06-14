package mongodb_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
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

func (m *TrainerTestSuite) TestShouldInsertTrainerWorkoutGroupWhenNotExist() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workoutGroup := testutil.GenerateTrainerWorkoutGroup(trainerUUID)
	// when:
	err := m.commandHandler.UpsertWorkoutGroup(ctx, workoutGroup)

	// then:
	assertions.Nil(err)

	actualSchedule, err := m.queryHandler.QueryWorkoutGroup(ctx, workoutGroup.UUID(), workoutGroup.TrainerUUID())
	assertions.Nil(err)
	assertions.Equal(workoutGroup, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldUpdateNameOfExistingTrainerWorkoutGroup() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workoutGroup := testutil.GenerateTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, workoutGroup)
	_ = workoutGroup.UpdateDesc("dummy2")

	// when:
	err := m.commandHandler.UpsertWorkoutGroup(ctx, workoutGroup)

	// then:
	assertions.Nil(err)

	actualSchedule, err := m.queryHandler.QueryWorkoutGroup(ctx, workoutGroup.UUID(), workoutGroup.TrainerUUID())
	assertions.Nil(err)
	assertions.Equal(workoutGroup, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldDeleteTrainerWorkoutGroupWithSuccess() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	workoutGroup := testutil.GenerateTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, workoutGroup)

	// when:
	err := m.commandHandler.DeleteWorkoutGroup(ctx, workoutGroup.UUID(), workoutGroup.TrainerUUID())

	// then:
	assertions.Nil(err)

	actualSchedule, err := m.queryHandler.QueryWorkoutGroup(ctx, workoutGroup.UUID(), workoutGroup.TrainerUUID())
	assertions.Nil(err)
	assertions.Empty(actualSchedule)
}

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroupWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const fakeUUID = "13dd31ee-e131-44e1-8d95-dd6317af81b7"
	const fakeTrainerUUID = "b71f6ed1-a982-47f2-a3fe-1d32cf3a132f"
	ctx := context.Background()
	workoutGroup := testutil.GenerateTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, workoutGroup)

	// when:
	err := m.commandHandler.DeleteWorkoutGroup(ctx, fakeUUID, fakeTrainerUUID)

	// then:
	assertions.Nil(err)

	actualSchedule, err := m.queryHandler.QueryWorkoutGroup(ctx, workoutGroup.UUID(), workoutGroup.TrainerUUID())
	assertions.Nil(err)
	assertions.Equal(workoutGroup, actualSchedule)
}

func (m *TrainerTestSuite) TestShouldDeleteAllTrainerWorkoutGroupsWithSuccess() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerWorkoutGroup(trainerUUID)
	second := testutil.GenerateTrainerWorkoutGroup(trainerUUID)

	_ = m.commandHandler.UpsertWorkoutGroup(ctx, first)
	_ = m.commandHandler.UpsertWorkoutGroup(ctx, second)

	// when:
	err := m.commandHandler.DeleteWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualSchedules, err := m.queryHandler.QueryWorkoutGroups(ctx, trainerUUID)
	assertions.Nil(err)
	assertions.Empty(actualSchedules)
}

func (m *TrainerTestSuite) TestShouldReturnEmptyWorkoutGroupsWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assertions := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	err := m.commandHandler.DeleteWorkoutGroups(ctx, trainerUUID)

	// then:
	assertions.Nil(err)

	actualSchedules, err := m.queryHandler.QueryWorkoutGroups(ctx, trainerUUID)
	assertions.Nil(err)
	assertions.Empty(actualSchedules)
}
