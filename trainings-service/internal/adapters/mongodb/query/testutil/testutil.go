package testutil

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/query"
	rm "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	DatabaseName   = "trainings_service_test"
	CollectionName = "trainings"
)

func CreateParticipantTrainingGroups(cli *mongo.Client, UUID string) ([]rm.ParticipantGroup, error) {
	writeModels, err := FindAllTrainingGroupsWithParticipant(cli, UUID)
	if err != nil {
		return nil, err
	}
	participantTrainingGroups := query.ConvertToParticipantGroups(writeModels...)
	return participantTrainingGroups, nil
}

func CreateTrainerGroup(cli *mongo.Client, groupUUID string) rm.TrainerGroup {
	writeModel, _ := FindTrainingGroup(cli, groupUUID)
	readModel := query.ConvertToQueryTrainerWorkoutGroup(writeModel)
	return readModel
}

func CreateAllTrainingGroups(cli *mongo.Client) ([]rm.TrainingGroup, error) {
	writeModels, err := FindAllTrainingGroups(cli)
	if err != nil {
		return nil, err
	}
	allTrainingGroups := query.ConvertToQueryTrainingGroups(writeModels...)
	return allTrainingGroups, nil
}

func FindAllTrainingGroupsWithParticipant(cli *mongo.Client, UUID string) ([]documents.TrainingGroupWriteModel, error) {
	db := cli.Database(DatabaseName)
	coll := db.Collection(CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	f := bson.M{"participants._id": UUID}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var docs []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func FindAllTrainingGroups(cli *mongo.Client) ([]documents.TrainingGroupWriteModel, error) {
	db := cli.Database(DatabaseName)
	coll := db.Collection(CollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	f := bson.D{}
	cur, err := coll.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var docs []documents.TrainingGroupWriteModel
	err = cur.All(ctx, &docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func FindTrainingGroup(cli *mongo.Client, uuid string) (documents.TrainingGroupWriteModel, error) {
	db := cli.Database(DatabaseName)
	coll := db.Collection(CollectionName)

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

func NewTestTrainingGroup(UUID string, trainer trainings.Trainer, date time.Time) trainings.TrainingGroup {
	t, err := trainings.NewTrainingGroup(UUID, "dummy name", "dummy desc", date, trainer)
	if err != nil {
		panic(err)
	}
	return *t
}

func NewTestStaticTime() time.Time {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	return ts
}

func NewTestTrainer(UUID, name string) trainings.Trainer {
	t, err := trainings.NewTrainer(UUID, name)
	if err != nil {
		panic(err)
	}
	return t
}

func NewTestParticipant(UUID, name string) trainings.Participant {
	p, err := trainings.NewParticipant(UUID, name)
	if err != nil {
		panic(err)
	}
	return p
}

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
