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

func (m *MongoDBTestSuite) TestShouldInsertTrainerWorkoutGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer, _ := trainings.NewTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	staticTime, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	expectedWorkoutGroup, _ := trainings.NewWorkoutGroup("76740131-ff8c-477b-895e-c9b80b08858c", "dummy name", "dummy desc", staticTime, trainer)

	// when:
	err := m.SUT.InsertTrainerWorkoutGroup(ctx, expectedWorkoutGroup)

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainerWorkoutGroupWriteModel(expectedWorkoutGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup, err := m.convertTrainerWorkoutGroupWriteModelToDomain(writeModel)
	assertions.Nil(err)
	assertions.Equal(*expectedWorkoutGroup, actualWorkoutDomainGroup)
}

func (m *MongoDBTestSuite) TestShouldQueryTrainerWorkoutGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer, _ := trainings.NewTrainer("914c6f44-8715-4ca2-9e35-0bae22ac52c3", "John Doe")
	staticTime, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	expectedWorkoutGroup, _ := trainings.NewWorkoutGroup("6b8e54fe-8727-463d-8ff3-7f1003eeee87", "dummy name", "dummy desc", staticTime, trainer)

	_ = m.SUT.InsertTrainerWorkoutGroup(ctx, expectedWorkoutGroup)

	// when:
	actualWorkoutGroup, err := m.SUT.QueryTrainerWorkoutGroup(ctx, expectedWorkoutGroup.UUID(), expectedWorkoutGroup.Trainer().UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(*expectedWorkoutGroup, actualWorkoutGroup)

	writeModel, err := m.findTrainerWorkoutGroupWriteModel(expectedWorkoutGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup, err := m.convertTrainerWorkoutGroupWriteModelToDomain(writeModel)
	assertions.Nil(err)
	assertions.Equal(*expectedWorkoutGroup, actualWorkoutDomainGroup)
}

func (m *MongoDBTestSuite) TestShouldReturnEmptyScheduledTrainerWorkoutGroupWhenNotScheduledBefore() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	const trainerUUID = "914c6f44-8715-4ca2-9e35-0bae22ac52c3"
	const groupUUID = "6b8e54fe-8727-463d-8ff3-7f1003eeee87"

	// when:
	actualWorkoutGroup, err := m.SUT.QueryTrainerWorkoutGroup(ctx, groupUUID, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Empty(actualWorkoutGroup)

	writeModel, err := m.findTrainerWorkoutGroupWriteModel(groupUUID)
	assertions.Equal(err, mongo.ErrNoDocuments)
	assertions.Empty(writeModel)
}

func (m *MongoDBTestSuite) TestShouldReturnTrainerWorkoutGroupQueryReadModelWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer, _ := trainings.NewTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	staticTime, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	group, _ := trainings.NewWorkoutGroup("76740131-ff8c-477b-895e-c9b80b08858c", "dummy name", "dummy desc", staticTime, trainer)

	_ = m.SUT.InsertTrainerWorkoutGroup(ctx, group)
	expectedReadModel := m.createQueryWorkoutGroupReadModel(group.UUID())

	// when:
	actualReadModel, err := m.SUT.TrainerWorkoutGroup(ctx, group.UUID(), group.Trainer().UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedReadModel, actualReadModel)
}

func (m *MongoDBTestSuite) TestShouldUpdateTrainerWorkoutGroupWithSuccess() {
	assertions := m.Assert()

	// given:
	ctx := context.Background()
	trainer, _ := trainings.NewTrainer("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")
	staticTime, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	workoutGroup, _ := trainings.NewWorkoutGroup("76740131-ff8c-477b-895e-c9b80b08858c", "dummy name", "dummy desc", staticTime, trainer)

	updatedWorkoutGroup := workoutGroup
	participant, _ := trainings.NewParticipant("#0a4e9c95-1e13-491a-b8ff-c0536b5f8dd6", "Jerry Smith")
	_ = updatedWorkoutGroup.AssignParticipant(participant)

	_ = m.SUT.InsertTrainerWorkoutGroup(ctx, workoutGroup)

	// when:
	err := m.SUT.UpdateTrainerWorkoutGroup(ctx, updatedWorkoutGroup)

	// then:
	assertions.Nil(err)

	writeModel, err := m.findTrainerWorkoutGroupWriteModel(updatedWorkoutGroup.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup, err := m.convertTrainerWorkoutGroupWriteModelToDomain(writeModel)
	assertions.Nil(err)
	assertions.Equal(*updatedWorkoutGroup, actualWorkoutDomainGroup)
}

func (m *MongoDBTestSuite) createQueryWorkoutGroupReadModel(groupUUID string) query.TrainerWorkoutGroup {
	writeModel, _ := m.findTrainerWorkoutGroupWriteModel(groupUUID)
	readModel := mongodb.UnmarshalToQueryTrainerWorkoutGroup(writeModel)
	return readModel
}

func (m *MongoDBTestSuite) findTrainerWorkoutGroupWriteModel(uuid string) (mongodb.WorkoutGroupWriteModel, error) {
	db := m.testCli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)
	f := bson.M{"_id": uuid}
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.QueryTimeout)
	defer cancel()
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return mongodb.WorkoutGroupWriteModel{}, res.Err()
	}

	var doc mongodb.WorkoutGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return mongodb.WorkoutGroupWriteModel{}, err
	}
	return doc, nil
}

func (m *MongoDBTestSuite) convertTrainerWorkoutGroupWriteModelToDomain(d mongodb.WorkoutGroupWriteModel) (trainings.WorkoutGroup, error) {
	var pp []trainings.DatabaseWorkoutGroupParticipant
	for _, p := range d.Participants {
		pp = append(pp, trainings.DatabaseWorkoutGroupParticipant{UUID: p.UUID, Name: p.Name})
	}
	g := trainings.UnmarshalWorkoutGroupFromDatabase(trainings.DatabaseWorkoutGroup{
		UUID:        d.UUID,
		Name:        d.Name,
		Description: d.Description,
		Limit:       d.Limit,
		Date:        d.Date,
		Trainer: trainings.DatabaseWorkoutGroupTrainer{
			UUID: d.Trainer.UUID,
			Name: d.Trainer.Name,
		},
		Participants: pp,
	})
	return g, nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMongoDBTestSuite_Integration(t *testing.T) {
	cfg := mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "trainings_service_test",
		Collection: "customer_schedules",
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
