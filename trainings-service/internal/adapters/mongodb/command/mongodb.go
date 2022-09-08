package command

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"
)

type Config struct {
	Database       string
	Collection     string
	CommandTimeout time.Duration
}

func ConvertToWriteModelParticipants(pp ...trainings.Participant) []documents.ParticipantWriteModel {
	var out []documents.ParticipantWriteModel
	for _, p := range pp {
		out = append(out, documents.ParticipantWriteModel{
			UUID: p.UUID(),
			Name: p.Name(),
		})
	}
	return out
}
