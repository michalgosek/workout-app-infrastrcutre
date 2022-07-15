package trainer

type WorkoutGroupWriteModel struct {
	UUID         string                              `bson:"_id"`
	TrainerUUID  string                              `bson:"trainer_uuid"`
	TrainerName  string                              `bson:"trainer_name"`
	Name         string                              `bson:"group_name"`
	Description  string                              `bson:"group_desc"`
	Limit        int                                 `bson:"limit"`
	Participants []WorkoutGroupParticipantWriteModel `bson:"participants"`
	Date         string                              `bson:"date"`
}

type WorkoutGroupParticipantWriteModel struct {
	UUID string `bson:"uuid"`
	Name string `bson:"name"`
}
