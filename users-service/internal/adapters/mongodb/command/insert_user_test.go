package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/command"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldInsertUserWithSuccess_Integration(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	user := testutil.NewTestUser("a6ae7d84-2938-4291-ae28-cb92ceba4f59", "John Doe")

	cli := testutil.NewTestMongoClient()
	SUT := command.NewInsertUserHandler(cli, command.Config{
		Database:       "users_service_test_db",
		Collection:     "users",
		CommandTimeout: 5 * time.Second,
	})

	defer func() {
		db := cli.Database("users_service_test_db")
		err := db.Drop(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// when:
	err := SUT.Do(ctx, &user)

	// then:
	assertions.Nil(err)

	writeModel, err := testutil.FindUser(cli, user.UUID())
	assertions.Nil(err)
	assertions.NotEmpty(writeModel)

	actualWorkoutDomainGroup := command.UnmarshalToUser(writeModel)
	assertions.Equal(user, actualWorkoutDomainGroup)
}
