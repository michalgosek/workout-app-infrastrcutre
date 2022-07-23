package testutil

import (
	"context"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/users-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"time"
)

func NewTestMongoClient() *mongo.Client {
	opts := options.Client()
	opts.ApplyURI("mongodb://localhost:27017")
	opts.SetConnectTimeout(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		panic(err)
	}
	err = cli.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return cli
}

func NewTestUser(UUID, role string) domain.User {
	n := rand.Int()
	email := fmt.Sprintf("email%d@email.com", n)
	name := fmt.Sprintf("test_user_%d", n)
	u, err := domain.NewUser(UUID, role, name, email)
	if err != nil {
		panic(err)
	}
	return *u
}

func FindUser(cli *mongo.Client, uuid string) (documents.UserWriteModel, error) {
	db := cli.Database("users_service_test_db")
	coll := db.Collection("users")

	f := bson.M{"_id": uuid}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return documents.UserWriteModel{}, res.Err()
	}

	var doc documents.UserWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return documents.UserWriteModel{}, err
	}
	return doc, nil
}
