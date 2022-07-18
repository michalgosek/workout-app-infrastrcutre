package mongodb

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func NewClient(addr string, d time.Duration) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(addr)
	opts.SetConnectTimeout(d)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cli, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}
	err = cli.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = cli.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func UnmarshalToTrainingGroup(d TrainingGroupWriteModel) trainings.TrainingGroup {
	var pp []trainings.DatabaseTrainingGroupParticipant
	for _, p := range d.Participants {
		pp = append(pp, trainings.DatabaseTrainingGroupParticipant{UUID: p.UUID, Name: p.Name})
	}
	g := trainings.UnmarshalTrainingGroupFromDatabase(trainings.DatabaseTrainingGroup{
		UUID:        d.UUID,
		Name:        d.Name,
		Description: d.Description,
		Limit:       d.Limit,
		Date:        d.Date,
		Trainer: trainings.DatabaseTrainingGroupTrainer{
			UUID: d.Trainer.UUID,
			Name: d.Trainer.Name,
		},
		Participants: pp,
	})
	return g
}

func UnmarshalToWriteModelParticipants(pp ...trainings.Participant) []ParticipantWriteModel {
	var out []ParticipantWriteModel
	for _, p := range pp {
		out = append(out, ParticipantWriteModel{
			UUID: p.UUID(),
			Name: p.Name(),
		})
	}
	return out
}

func UnmarshalToQueryTrainerWorkoutGroup(d TrainingGroupWriteModel) query.TrainerWorkoutGroup {
	var pp []query.Participant
	for _, p := range d.Participants {
		pp = append(pp, query.Participant{
			Name: p.Name,
			UUID: p.UUID,
		})
	}
	g := query.TrainerWorkoutGroup{
		UUID:         d.UUID,
		Name:         d.Name,
		Description:  d.Description,
		Date:         d.Date,
		Limit:        d.Limit,
		Participants: pp,
	}
	return g
}

func UnmarshalToQueryTrainerWorkoutGroups(dd ...TrainingGroupWriteModel) []query.TrainerWorkoutGroup {
	var out []query.TrainerWorkoutGroup
	for _, d := range dd {
		g := UnmarshalToQueryTrainerWorkoutGroup(d)
		out = append(out, g)
	}
	return out
}
