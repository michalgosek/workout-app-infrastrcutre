package query

type TrainerWorkoutGroupParticipant struct {
	UUID string
	Name string
}

type TrainerWorkoutGroup struct {
	TrainerUUID  string
	TrainerName  string
	GroupUUID    string
	GroupDesc    string
	GroupName    string
	Date         string
	Participants []TrainerWorkoutGroupParticipant
}
