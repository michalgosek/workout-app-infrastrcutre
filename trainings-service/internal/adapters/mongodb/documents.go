package mongodb

import "time"

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
