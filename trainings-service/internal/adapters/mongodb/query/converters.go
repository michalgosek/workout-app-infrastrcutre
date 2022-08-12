package query

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
)

func ConvertToQueryTrainerWorkoutGroup(d documents.TrainingGroupWriteModel) query.TrainerGroup {
	var pp []query.Participant
	for _, p := range d.Participants {
		pp = append(pp, query.Participant{
			Name: p.Name,
			UUID: p.UUID,
		})
	}
	g := query.TrainerGroup{
		UUID:         d.UUID,
		Name:         d.Name,
		Description:  d.Description,
		Date:         d.Date.Format(query.UIFormat),
		Limit:        d.Limit,
		Participants: pp,
	}
	return g
}

func ConvertToTrainerWorkoutGroups(dd ...documents.TrainingGroupWriteModel) []query.TrainerGroup {
	var out []query.TrainerGroup
	for _, d := range dd {
		g := ConvertToQueryTrainerWorkoutGroup(d)
		out = append(out, g)
	}
	return out
}

func ConvertToParticipantGroups(dd ...documents.TrainingGroupWriteModel) []query.ParticipantGroup {
	var out []query.ParticipantGroup
	for _, d := range dd {
		out = append(out, query.ParticipantGroup{
			UUID:        d.UUID,
			TrainerUUID: d.Trainer.UUID,
			TrainerName: d.Trainer.Name,
			Name:        d.Name,
			Description: d.Description,
			Date:        d.Date.Format(query.UIFormat),
		})
	}
	return out
}

func ConvertToQueryTrainingGroups(dd ...documents.TrainingGroupWriteModel) []query.TrainingGroup {
	var out []query.TrainingGroup
	for _, d := range dd {
		var participants []query.Participant
		for _, p := range d.Participants {
			participants = append(participants, query.Participant{
				Name: p.Name,
				UUID: p.UUID,
			})
		}
		out = append(out, query.TrainingGroup{
			UUID:         d.UUID,
			TrainerUUID:  d.Trainer.UUID,
			TrainerName:  d.Trainer.Name,
			Name:         d.Name,
			Description:  d.Description,
			Date:         d.Date.Format(query.UIFormat),
			Limit:        d.Limit,
			Participants: participants,
		})
	}
	return out
}
