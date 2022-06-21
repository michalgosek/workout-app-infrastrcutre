package trainer

type WorkoutGroupDocument struct {
	UUID            string                    `bson:"_id"`
	TrainerUUID     string                    `bson:"trainer_uuid"`
	TrainerName     string                    `bson:"trainer_name"`
	WorkoutName     string                    `bson:"workout_name"`
	WorkoutDesc     string                    `bson:"workout_desc"`
	Limit           int                       `bson:"limit"`
	CustomerDetails []CustomerDetailsDocument `bson:"customer_details"`
	Date            string                    `bson:"date"`
}

type CustomerDetailsDocument struct {
	UUID string `bson:"uuid"`
	Name string `bson:"name"`
}
