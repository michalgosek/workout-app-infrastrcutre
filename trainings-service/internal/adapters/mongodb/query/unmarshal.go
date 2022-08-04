package query

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/mongodb/documents"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"
)

func UnmarshalToQueryTrainerWorkoutGroup(d documents.TrainingGroupWriteModel) query.TrainerWorkoutGroup {
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

func UnmarshalToQueryTrainerWorkoutGroups(dd ...documents.TrainingGroupWriteModel) []query.TrainerWorkoutGroup {
	var out []query.TrainerWorkoutGroup
	for _, d := range dd {
		g := UnmarshalToQueryTrainerWorkoutGroup(d)
		out = append(out, g)
	}
	return out
}

func UnmarshalToQueryTrainingGroups(dd ...documents.TrainingGroupWriteModel) []query.TrainingWorkoutGroup {
	var out []query.TrainingWorkoutGroup
	for _, d := range dd {
		out = append(out, query.TrainingWorkoutGroup{
			UUID:         d.UUID,
			TrainerUUID:  d.Trainer.UUID,
			TrainerName:  d.Trainer.Name,
			Name:         d.Name,
			Description:  d.Description,
			Date:         d.Date,
			Limit:        d.Limit,
			Participants: len(d.Participants),
		})
	}
	return out
}
