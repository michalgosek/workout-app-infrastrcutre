package mongodb

type TrainerScheduleDocument struct {
	UUID          string   `bson:"_id"`
	TrainerUUID   string   `bson:"trainer_uuid"`
	Limit         int      `bson:"limit"`
	CustomerUUIDs []string `bson:"customer_uuids"`
	Name          string   `bson:"name"`
	Desc          string   `bson:"desc"`
	Date          string   `bson:"date"`
}
