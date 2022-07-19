package documents

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type ParticipantWriteModel struct {
	UUID string `bson:"_id"`
	Name string `bson:"name"`
}

type TrainerWriteModel struct {
	UUID           string                    `bson:"_id"`
	Name           string                    `bson:"name"`
	TrainingGroups []TrainingGroupWriteModel `bson:"training_groups,omitempty"`
}

type TrainingGroupWriteModel struct {
	UUID         string                  `bson:"_id"`
	Name         string                  `bson:"name"`
	Description  string                  `bson:"description"`
	Date         time.Time               `bson:"date"`
	Trainer      TrainerWriteModel       `bson:"trainer"`
	Limit        int                     `bson:"limit"`
	Participants []ParticipantWriteModel `bson:"participants"`
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
