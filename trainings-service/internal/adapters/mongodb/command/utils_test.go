package command_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func findTrainingGroup(cli *mongo.Client, uuid string) (documents.TrainingGroupWriteModel, error) {
	db := cli.Database("insert_training_db")
	coll := db.Collection("trainings")

	f := bson.M{"_id": uuid}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := coll.FindOne(ctx, f)
	if res.Err() != nil {
		return documents.TrainingGroupWriteModel{}, res.Err()
	}

	var doc documents.TrainingGroupWriteModel
	err := res.Decode(&doc)
	if err != nil {
		return documents.TrainingGroupWriteModel{}, err
	}
	return doc, nil
}

func newTestParticipant(UUID string) trainings.Participant {
	p, err := trainings.NewParticipant(UUID, "Jerry Smith")
	if err != nil {
		panic(err)
	}
	return p
}

func newTestTrainingGroup(UUID string, trainer trainings.Trainer, date time.Time) trainings.TrainingGroup {
	t, err := trainings.NewTrainingGroup(UUID, "dummy name", "dummy desc", date, trainer)
	if err != nil {
		panic(err)
	}
	return *t
}

func newTestStaticTime() time.Time {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	return ts
}

func newTestTrainer(UUID, name string) trainings.Trainer {
	t, err := trainings.NewTrainer(UUID, name)
	if err != nil {
		panic(err)
	}
	return t
}

func newTestMongoClient() *mongo.Client {
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
