package mongodb

import "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/query"

func UnmarshalToQueryTrainerWorkoutGroup(d WorkoutGroupWriteModel) query.TrainerWorkoutGroup {
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

func UnmarshalToQueryTrainerWorkoutGroups(dd ...WorkoutGroupWriteModel) []query.TrainerWorkoutGroup {
	var out []query.TrainerWorkoutGroup
	for _, d := range dd {
		g := UnmarshalToQueryTrainerWorkoutGroup(d)
		out = append(out, g)
	}
	return out
}
