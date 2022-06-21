package query

type CustomerData struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type WorkoutGroupDetails struct {
	TrainerUUID string         `json:"trainer_uuid"`
	TrainerName string         `json:"trainer_name"`
	GroupUUID   string         `json:"group_uuid"`
	GroupDesc   string         `json:"group_desc"`
	GroupName   string         `json:"group_name"`
	Customers   []CustomerData `json:"customers"`
	Date        string         `json:"date"`
}
