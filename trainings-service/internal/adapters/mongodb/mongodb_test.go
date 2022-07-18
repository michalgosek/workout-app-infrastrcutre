package mongodb_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type MongoDBTestSuite struct {
	suite.Suite
	testCli *mongo.Client
	cfg     mongodb.Config
	SUT     *mongodb.Repository
}

func (m *MongoDBTestSuite) SetupTest() {
	t := m.Suite.T()
	db := m.testCli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.CommandTimeout)
	defer cancel()

	err := coll.Drop(ctx)
	if err != nil {
		t.Fatalf("mongo cli collection %s drop failed: %s", m.cfg.Collection, err)
	}
}

func (m *MongoDBTestSuite) TearDownSuite() {
	t := m.Suite.T()
	db := m.testCli.Database(m.cfg.Database)

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.CommandTimeout)
	defer cancel()
	err := db.Drop(ctx)
	if err != nil {
		t.Fatalf("mongo cli database %s drop failed: %s", m.cfg.Database, err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), m.cfg.Timeouts.CommandTimeout)
	defer cancel()
	err = m.testCli.Disconnect(ctx)
	if err != nil {
		t.Fatalf("mongo cli disconnect failed: %s", err)
	}
}

func (m *MongoDBTestSuite) TestShouldInsertTrainingGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	// when:
	err := m.SUT.InsertTrainingGroup(ctx, &trainingGroup)

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainingGroup(trainingGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup := mongodb.UnmarshalToTrainingGroup(writeModel)
	assertions.Equal(trainingGroup, actualWorkoutDomainGroup)
}

func (m *MongoDBTestSuite) TestShouldQueryTrainingGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("cd2c48da-ec19-4c32-8846-9aa85d1eeff3", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("7f340572-654f-4280-a0ad-b66bb70bd1ac", trainer, date)

	_ = m.SUT.InsertTrainingGroup(ctx, &trainingGroup)

	// when:
	actualWorkoutGroup, err := m.SUT.QueryTrainingGroup(ctx, trainingGroup.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(trainingGroup, actualWorkoutGroup)
}

func (m *MongoDBTestSuite) TestShouldReturnEmptyTrainingGroupReadModelWhenNonExist() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	const groupUUID = "6b8e54fe-8727-463d-8ff3-7f1003eeee87"

	// when:
	actualQueryTrainingGroup, err := m.SUT.QueryTrainingGroup(ctx, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualQueryTrainingGroup)
}

func (m *MongoDBTestSuite) TestShouldReturnTrainingGroupReadModelWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	_ = m.SUT.InsertTrainingGroup(ctx, &trainingGroup)
	expectedReadModel := m.createTrainingGroupReadModel(trainingGroup.UUID())

	// when:
	actualReadModel, err := m.SUT.TrainingGroup(ctx, trainingGroup.UUID(), trainingGroup.Trainer().UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedReadModel, actualReadModel)
}

func (m *MongoDBTestSuite) TestShouldUpdateTrainingGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	_ = m.SUT.InsertTrainingGroup(ctx, &trainingGroup)
	_ = trainingGroup.AssignParticipant(newTestParticipant("c6975a21-a098-4c94-a7de-de01a731b57a"))

	// when:
	err := m.SUT.UpdateTrainingGroup(ctx, &trainingGroup)

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainingGroup(trainingGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup := mongodb.UnmarshalToTrainingGroup(writeModel)
	assertions.Equal(trainingGroup, actualWorkoutDomainGroup)
}

func (m *MongoDBTestSuite) TestShouldDeleteTrainingGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	trainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)

	_ = m.SUT.InsertTrainingGroup(ctx, &trainingGroup)

	// when:
	err := m.SUT.DeleteTrainingGroup(ctx, trainingGroup.UUID(), trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainingGroup(trainingGroup.UUID())
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}

func (m *MongoDBTestSuite) TestShouldDeleteTrainingGroupsWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer := newTestTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	date := newTestStaticTime()
	firstTrainingGroup := newTestTrainingGroup("76740131-ff8c-477b-895e-c9b80b08858c", trainer, date)
	secondTrainingGroup := newTestTrainingGroup("76275b9a-2c9c-4d96-adb7-ab44ad2ab8ee", trainer, date)

	_ = m.SUT.InsertTrainingGroup(ctx, &firstTrainingGroup)
	_ = m.SUT.InsertTrainingGroup(ctx, &secondTrainingGroup)

	// when:
	err := m.SUT.DeleteTrainingGroups(ctx, trainer.UUID())

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainingGroups(trainer.UUID())
	assertions.Nil(err)
	assertions.Empty(writeModel)
}

func (m *MongoDBTestSuite) createTrainingGroupReadModel(groupUUID string) query.TrainerWorkoutGroup {
	writeModel, _ := m.findTrainingGroup(groupUUID)
	readModel := mongodb.UnmarshalToQueryTrainerWorkoutGroup(writeModel)
	return readModel
}

func (m *MongoDBTestSuite) findTrainingGroups(trainerUUID string) ([]mongodb.TrainingGroupWriteModel, error) {
	db := m.testCli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"trainer._id": trainerUUID}
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.QueryTimeout)
	defer cancel()

	curr, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	var dst []mongodb.TrainingGroupWriteModel
	err = curr.All(ctx, &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (m *MongoDBTestSuite) findTrainingGroup(uuid string) (mongodb.TrainingGroupWriteModel, error) {
	db := m.testCli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": uuid}
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.QueryTimeout)
	defer cancel()
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return mongodb.TrainingGroupWriteModel{}, res.Err()
	}

	var doc mongodb.TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return mongodb.TrainingGroupWriteModel{}, err
	}
	return doc, nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMongoDBTestSuite_Integration(t *testing.T) {
	cfg := mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "trainings_service_test",
		Collection: "trainings",
		Timeouts: mongodb.Timeouts{
			CommandTimeout:    10 * time.Second,
			QueryTimeout:      10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
	}

	cli, err := mongodb.NewClient(cfg.Addr, cfg.Timeouts.ConnectionTimeout)
	if err != nil {
		t.Fatalf("creating mongo client failed: %s", err)
	}
	SUT, err := mongodb.NewRepository(cfg)
	if err != nil {
		t.Fatalf("creating mongo repository  failed: %s", err)
	}
	ts := MongoDBTestSuite{
		cfg:     cfg,
		testCli: cli,
		SUT:     SUT,
	}
	suite.Run(t, &ts)
}

func newTestTrainer(UUID, name string) trainings.Trainer {
	t, err := trainings.NewTrainer(UUID, name)
	if err != nil {
		panic(err)
	}
	return t
}

func newTestParticipant(UUID string) trainings.Participant {
	p, err := trainings.NewParticipant("0a4e9c95-1e13-491a-b8ff-c0536b5f8dd6", "Jerry Smith")
	if err != nil {
		panic(err)
	}
	return p
}

func newTestStaticTime() time.Time {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	return ts
}

func newTestTrainingGroup(UUID string, trainer trainings.Trainer, date time.Time) trainings.TrainingGroup {
	t, err := trainings.NewTrainingGroup(UUID, "dummy name", "dummy desc", date, trainer)
	if err != nil {
		panic(err)
	}
	return *t
}
