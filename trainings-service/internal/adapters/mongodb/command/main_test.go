package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	cli := testutil.NewTestMongoClient()
	db := cli.Database(testutil.DatabaseName)
	err := db.Drop(ctx)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}
