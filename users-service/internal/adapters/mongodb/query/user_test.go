package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/query"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/testutil"
	rm "github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/application/query"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestShouldReturnUserReadModelWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	user := testutil.NewTestUser("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")

	cli := testutil.NewTestMongoClient()
	insertTrainingHandler := command.NewInsertUserHandler(cli, command.Config{
		Database:       "users_service_test_db",
		Collection:     "users",
		CommandTimeout: 5 * time.Second,
	})
	_ = insertTrainingHandler.Do(ctx, &user)

	SUT := query.NewUserHandler(cli, query.Config{
		Database:     "users_service_test_db",
		Collection:   "users",
		QueryTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database("users_service_test_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	expectedReadModel := createUserReadModel(cli, user.UUID())

	// when:
	training, err := SUT.User(ctx, user.UUID())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedReadModel, training)
}

func createUserReadModel(cli *mongo.Client, groupUUID string) rm.User {
	writeModel, _ := testutil.FindUser(cli, groupUUID)
	readModel := query.UnmarshalToQueryUser(writeModel)
	return readModel
}
