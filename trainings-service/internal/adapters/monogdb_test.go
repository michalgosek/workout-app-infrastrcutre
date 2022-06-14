package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMongoDBSuite_Integration(t *testing.T) {
	suite.Run(t, new(MongoDBSuite))
}

type MongoDBSuite struct {
	suite.Suite
	SUT *adapters.MongoDB
}

func (m *MongoDBSuite) BeforeTest(mName, testName string) {
	t := m.T()
	cfg := adapters.MongoDBConfig{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service_test",
		Collection:        "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Format:            "2006-01-02 15:04",
	}

	SUT, err := adapters.NewMongoDB(cfg)
	if err != nil {
		t.Fatalf("creating SUT failed: %v", err)
	}

	err = SUT.Ping(context.Background())
	if err != nil {
		t.Fatalf("SUT ping request failed: %v", err)
	}

	m.SUT = SUT
	t.Logf("SUT connection established with mongodb node on addr: %s", cfg.Addr)
}

func (m *MongoDBSuite) AfterTest(suiteName, testName string) {
	t := m.T()
	ctx := context.Background()
	err := m.SUT.DropCollection(ctx)
	if err != nil {
		t.Fatalf("SUT dropping collection failed: %v", err)
	}
	t.Log("SUT cli dropped collection successfully")
}

func (m *MongoDBSuite) TearDownSuite() {
	t := m.T()
	ctx := context.Background()

	err := m.SUT.DropDatabase(ctx)
	if err != nil {
		t.Logf("SUT cli dropping db failed: %v", err)
	}
	t.Logf("SUT cli dropped db successfully")

	err = m.SUT.Disconnect(ctx)
	if err != nil {
		t.Logf("SUT cli disconnect failed: %v", err)
	}
	t.Logf("SUT cli disconnect successfully")
}

func (m *MongoDBSuite) TestShouldInsertTrainerScheduleWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	// when:
	err := m.SUT.UpsertTrainerSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *MongoDBSuite) TestShouldUpdateNameOfExisitingSchedule() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertTrainerSchedule(ctx, expectedSchedule)
	expectedSchedule.UpdateDesc("dummy2")

	// when:
	err := m.SUT.UpsertTrainerSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *MongoDBSuite) TestShouldCancelTrainerScheduleWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertTrainerSchedule(ctx, expectedSchedule)

	// when:
	err := m.SUT.CancelTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Empty(actualSchedule)
}

func (m *MongoDBSuite) TestShouldReturnEmptyScheduleWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const fakeUUID = "13dd31ee-e131-44e1-8d95-dd6317af81b7"
	const fakeTrainerUUID = "b71f6ed1-a982-47f2-a3fe-1d32cf3a132f"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertTrainerSchedule(ctx, expectedSchedule)

	// when:
	err := m.SUT.CancelTrainerSchedule(ctx, fakeUUID, fakeTrainerUUID)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QueryTrainerSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func (m *MongoDBSuite) TestShouldReturnTwoSchedulesWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)
	expectedSchedules := []trainer.TrainerSchedule{first, second}

	m.SUT.UpsertTrainerSchedule(ctx, first)
	m.SUT.UpsertTrainerSchedule(ctx, second)

	// when:
	actualSchedules, err := m.SUT.QueryTrainerSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSchedules, actualSchedules)
}

func (m *MongoDBSuite) TestShouldReturnEmptySchedulesWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualSchedules, err := m.SUT.QueryTrainerSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

func (m *MongoDBSuite) TestShouldCancelAllUpsertedSchedules() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertTrainerSchedule(ctx, first)
	m.SUT.UpsertTrainerSchedule(ctx, second)

	// when:
	err := m.SUT.CancelTrainerSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.SUT.QueryTrainerSchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

func (m *MongoDBSuite) TestShouldReturnEmptySchedulesWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	err := m.SUT.CancelTrainerSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.SUT.QueryTrainerSchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}
