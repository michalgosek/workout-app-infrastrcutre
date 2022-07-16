package mongodb

import "time"

type ParticipantWriteModel struct {
	UUID string `bson:"_id"`
	Name string `bson:"name"`
}

type TrainerWorkoutGroupsWriteModel struct {
	UUID          string                   `bson:"_id"`
	Name          string                   `bson:"name"`
	WorkoutGroups []WorkoutGroupWriteModel `bson:"workout_groups,omitempty"`
}

type WorkoutGroupWriteModel struct {
	UUID         string                         `bson:"_id"`
	Name         string                         `bson:"name"`
	Description  string                         `bson:"description"`
	Date         time.Time                      `bson:"date"`
	Trainer      TrainerWorkoutGroupsWriteModel `bson:"trainer"`
	Limit        int                            `bson:"limit"`
	Participants []ParticipantWriteModel        `bson:"participants"`
}
