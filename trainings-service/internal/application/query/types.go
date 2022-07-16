package query

import "time"

type Participant struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type TrainerWorkoutGroup struct {
	UUID         string        `bson:"_id"`
	Name         string        `bson:"name"`
	Description  string        `bson:"description"`
	Date         time.Time     `bson:"date"`
	Limit        int           `bson:"limit"`
	Participants []Participant `bson:"participants"`
}
