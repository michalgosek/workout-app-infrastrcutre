package ports

type User struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type TrainingGroupPost struct {
	User      User   `json:"user"`
	GroupName string `json:"group_name"`
	GroupDesc string `json:"group_desc"`
	Date      string `json:"date"`
}

type UpdateTrainingGroupPost struct {
	GroupName string `json:"group_name"`
	GroupDesc string `json:"group_desc"`
	Date      string `json:"date"`
}
