package trainer

type WorkoutGroupDocument struct {
	UUID          string   `bson:"_id"`
	TrainerUUID   string   `bson:"trainer_uuid"`
	TrainerName   string   `bson:"trainer_name"`
	WorkoutName   string   `bson:"workout_name"`
	WorkoutDesc   string   `bson:"workout_desc"`
	Limit         int      `bson:"limit"`
	CustomerUUIDs []string `bson:"customer_uuids"`
	Date          string   `bson:"date"`
}
