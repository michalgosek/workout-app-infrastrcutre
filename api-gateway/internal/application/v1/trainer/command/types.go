package command

type PlanTraining struct {
	UserUUID  string
	GroupName string
	GroupDesc string
	Date      string
}

type User struct {
	UUID string
	Name string
	Role string
}

type PlanTrainingCommand struct {
	User      User
	GroupName string
	GroupDesc string
	Date      string
}
