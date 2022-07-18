package mongodb_test

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
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

func (m *MongoDBTestSuite) TestShouldInsertUserWithSuccess() {
	assertions := m.Assert()

	// given:
	const uuid = "e209ee9e-7fe4-4d45-b7a1-93622a1a74b7"
	ctx := context.Background()
	user := newTestUser(uuid, "trainer")

	// when:
	err := m.SUT.InsertUser(ctx, &user)

	// then:
	assertions.Nil(err)

	actualInsertedUser, err := m.findUser(uuid)
	assertions.Nil(err)
	assertions.Equal(user, actualInsertedUser)
}

func (m *MongoDBTestSuite) TestShouldQueryInsertedUserWithSuccess() {
	assertions := m.Assert()

	// given:
	const uuid = "a0fac876-784d-47ef-8508-1bad6841c8cf"
	ctx := context.Background()

	user := newTestUser(uuid, "trainer")
	_ = m.SUT.InsertUser(ctx, &user)

	// when:
	actualUser, err := m.SUT.QueryUser(ctx, uuid)

	// then
	assertions.Nil(err)
	assertions.Equal(user, actualUser)
}

func (m *MongoDBTestSuite) findUser(uuid string) (domain.User, error) {
	db := m.testCli.Database(m.cfg.Database)
	coll := db.Collection(m.cfg.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeouts.ConnectionTimeout)
	defer cancel()
	f := bson.M{"_id": uuid}
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return domain.User{}, nil
	}

	var dst mongodb.UserWriteModel
	err := res.Decode(&dst)
	if err != nil {
		return domain.User{}, nil
	}
	u := domain.UnmarshalUserFromDatabase(domain.DatabaseUser{
		UUID:           dst.UUID,
		Active:         dst.Active,
		Role:           dst.Role,
		Name:           dst.Name,
		LastActiveDate: dst.LastActiveDate,
	})
	return u, nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMongoDBTestSuite_Integration(t *testing.T) {
	cfg := mongodb.Config{
		Addr:       "mongodb://localhost:27017",
		Database:   "users_service_test",
		Collection: "users",
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

func newTestUser(UUID, role string) domain.User {
	n := rand.Int()
	name := fmt.Sprintf("test_user_%d", n)
	u, err := domain.NewUser(UUID, role, name)
	if err != nil {
		panic(err)
	}
	return *u
}
