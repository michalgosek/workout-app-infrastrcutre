package customer

type WorkoutDayDocument struct {
	UUID         string `bson:"_id"`
	CustomerName string `bson:"customer_name"`
	CustomerUUID string `bson:"customer_uuid"`
	GroupUUID    string `bson:"trainer_workout_group_uuid"`
	Date         string `bson:"date"`
}
