package mongodb

type TrainerWorkoutGroupDocument struct {
	UUID          string   `bson:"_id"`
	TrainerUUID   string   `bson:"trainer_uuid"`
	Limit         int      `bson:"limit"`
	CustomerUUIDs []string `bson:"customer_uuids"`
	Name          string   `bson:"name"`
	Desc          string   `bson:"desc"`
	Date          string   `bson:"date"`
}

type CustomerWorkoutDocument struct {
	UUID                    string `bson:"_id"`
	CustomerUUID            string `bson:"customer_uuid"`
	TrainerWorkoutGroupUUID string `bson:"trainer_workout_group_uuid"`
	Date                    string `bson:"date"`
}
