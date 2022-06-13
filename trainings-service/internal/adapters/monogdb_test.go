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

type MongoDBTrainerSchedulesSuite struct {
	suite.Suite
	SUT *adapters.TrainerSchedulesMongoDB
}

func (m *MongoDBTrainerSchedulesSuite) BeforeTest(mName, testName string) {
	t := m.T()
	cfg := adapters.MongoDBConfig{
		Addr:              "mongodb://localhost:27017",
		Database:          "trainings_service_test",
		Collection:        "trainer_schedules",
		CommandTimeout:    10 * time.Second,
		QueryTimeout:      10 * time.Second,
		ConnectionTimeout: 10 * time.Second,
	}

	SUT, err := adapters.NewTrainerSchedulesMongoDB(cfg)
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

func (m *MongoDBTrainerSchedulesSuite) AfterTest(mName, testName string) {
	t := m.T()
	ctx := context.Background()
	err := m.SUT.DropCollection(ctx)
	if err != nil {
		t.Fatalf("SUT dropping collection failed: %v", err)
	}
	t.Log("SUT cli dropped collection successfully")

	err = m.SUT.Disconnect(ctx)
	if err != nil {
		t.Logf("SUT cli disconnect failed: %v", err)
	}
	t.Logf("SUT cli disconnect successfully")
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldInsertTrainerScheduleWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	// when:
	err := m.SUT.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QuerySchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	AssertTrainerSchedule(t, expectedSchedule, actualSchedule)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldUpdateNameOfExisitingSchedule() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertSchedule(ctx, expectedSchedule)
	expectedSchedule.UpdateDesc("dummy2")

	// when:
	err := m.SUT.UpsertSchedule(ctx, expectedSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QuerySchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	AssertTrainerSchedule(t, expectedSchedule, actualSchedule)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldCancelTrainerScheduleWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertSchedule(ctx, expectedSchedule)

	// when:
	err := m.SUT.CancelSchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QuerySchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	assert.Empty(actualSchedule)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldReturnEmptyScheduleWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	const fakeUUID = "13dd31ee-e131-44e1-8d95-dd6317af81b7"
	const fakeTrainerUUID = "b71f6ed1-a982-47f2-a3fe-1d32cf3a132f"
	ctx := context.Background()
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertSchedule(ctx, expectedSchedule)

	// when:
	err := m.SUT.CancelSchedule(ctx, fakeUUID, fakeTrainerUUID)

	// then:
	assert.Nil(err)

	actualSchedule, err := m.SUT.QuerySchedule(ctx, expectedSchedule.UUID(), expectedSchedule.TrainerUUID())
	assert.Nil(err)
	AssertTrainerSchedule(t, expectedSchedule, actualSchedule)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldReturnTwoSchedulesWithSuccess() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)
	expectedSchedules := []trainer.TrainerSchedule{first, second}

	m.SUT.UpsertSchedule(ctx, first)
	m.SUT.UpsertSchedule(ctx, second)

	// when:
	actualSchedules, err := m.SUT.QuerySchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	AssertTrainerSchedules(t, expectedSchedules, actualSchedules)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldReturnEmptySchedulesWhenNotExist() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	actualSchedules, err := m.SUT.QuerySchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldCancelAllUpsertedSchedules() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()
	first := testutil.GenerateTrainerSchedule(trainerUUID)
	second := testutil.GenerateTrainerSchedule(trainerUUID)

	m.SUT.UpsertSchedule(ctx, first)
	m.SUT.UpsertSchedule(ctx, second)

	// when:
	err := m.SUT.CancelSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.SUT.QuerySchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

func (m *MongoDBTrainerSchedulesSuite) TestShouldReturnEmptySchedulesWhenNoExistForSpecifiedTrainer() {
	t := m.T()
	assert := assert.New(t)

	// given:
	const trainerUUID = "f1741a08-39d7-465d-adc9-a63cf058b409"
	ctx := context.Background()

	// when:
	err := m.SUT.CancelSchedules(ctx, trainerUUID)

	// then:
	assert.Nil(err)

	actualSchedules, err := m.SUT.QuerySchedules(ctx, trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedules)
}

// In order for 'go test' to run this m, we need to create
// a normal test function and pass our m to m.Run
func TestMongoDBTrainerSchedulesSuite_Integration(t *testing.T) {
	suite.Run(t, new(MongoDBTrainerSchedulesSuite))
}

func AssertTrainerSchedule(t *testing.T, expected, actual trainer.TrainerSchedule) {
	assert := assert.New(t)
	assert.Equal(expected.UUID(), actual.UUID())
	assert.Equal(expected.CustomerUUIDs(), actual.CustomerUUIDs())
	assert.Equal(expected.TrainerUUID(), actual.TrainerUUID())
	assert.Equal(expected.Name(), actual.Name())
	assert.Equal(expected.Limit(), actual.Limit())
	assert.Equal(expected.Desc(), actual.Desc())
}

func AssertTrainerSchedules(t *testing.T, expected, actual []trainer.TrainerSchedule) {
	assert := assert.New(t)
	assert.Len(actual, len(expected))
	for i := 0; i < len(expected); i++ {
		AssertTrainerSchedule(t, expected[i], actual[i])
	}
}
