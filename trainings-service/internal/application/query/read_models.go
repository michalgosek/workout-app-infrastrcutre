package query

import "time"

type Participant struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type TrainerWorkoutGroup struct {
	UUID         string        `json:"uuid"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Date         time.Time     `json:"date"`
	Limit        int           `json:"limit"`
	Participants []Participant `json:"participants"`
}
