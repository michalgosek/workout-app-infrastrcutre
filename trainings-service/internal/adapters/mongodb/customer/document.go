package customer

type WorkoutDocument struct {
	UUID                    string `bson:"_id"`
	CustomerUUID            string `bson:"customer_uuid"`
	TrainerWorkoutGroupUUID string `bson:"trainer_workout_group_uuid"`
	Date                    string `bson:"date"`
}
