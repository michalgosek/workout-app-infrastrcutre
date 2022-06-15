package http

type TrainerWorkoutGroupHTTPRequest struct {
	TrainerUUID string `json:"trainer_uuid"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Date        string `json:"date"`
}

type TrainerWorkoutGroupHTTPResponse struct {
	UUID          string   `json:"workout_group_uuid"`
	CustomerUUIDs []string `json:"customer_uuids"`
	Date          string   `json:"date"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Limit         int      `json:"limit"`
}

type TrainerWorkoutGroupsHTTPResponse struct {
	TrainerUUID   string                            `json:"trainer_uuid"`
	WorkoutGroups []TrainerWorkoutGroupHTTPResponse `json:"workout_groups"`
}

type DeleteTrainerWorkoutGroupHTTPResponse struct {
	UUID string `json:"workout_group_uuid"`
}

type DeleteTrainerWorkoutGroupsHTTPResponse struct {
	TrainerUUID string `json:"trainer_uuid_uuid"`
}
