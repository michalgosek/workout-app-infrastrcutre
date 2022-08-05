package query

type Participant struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type TrainerGroup struct {
	UUID         string        `json:"uuid"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Date         string        `json:"date"`
	Limit        int           `json:"limit"`
	Participants []Participant `json:"participants"`
}

type TrainingGroup struct {
	UUID         string `json:"uuid"`
	TrainerName  string `json:"trainer_name"`
	TrainerUUID  string `json:"trainer_uuid"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Date         string `json:"date"`
	Limit        int    `json:"limit"`
	Participants int    `json:"participants"`
}

type ParticipantGroup struct {
	UUID        string `json:"uuid"`
	TrainerName string `json:"trainer_name"`
	TrainerUUID string `json:"trainer_uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
