package trainer

import (
	"fmt"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type WorkoutGroupDocument struct {
	UUID          string   `bson:"_id"`
	TrainerUUID   string   `bson:"trainer_uuid"`
	Limit         int      `bson:"limit"`
	CustomerUUIDs []string `bson:"customer_uuids"`
	Name          string   `bson:"name"`
	Desc          string   `bson:"desc"`
	Date          string   `bson:"date"`
}

func convertToDomainWorkoutGroups(format string, docs ...WorkoutGroupDocument) ([]trainer.WorkoutGroup, error) {
	var workouts []trainer.WorkoutGroup
	for _, d := range docs {
		date, err := time.Parse(format, d.Date)
		if err != nil {
			return nil, fmt.Errorf("parsing date value from document failed: %v", err)
		}
		workout, err := trainer.UnmarshalFromDatabase(d.UUID, d.TrainerUUID, d.Name, d.Desc, d.CustomerUUIDs, date, d.Limit)
		if err != nil {
			return nil, fmt.Errorf("unmarshal from database failed: %v", err)
		}
		workouts = append(workouts, workout)
	}
	return workouts, nil
}
