package ports

type User struct {
	Name string `json:"name"`
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
