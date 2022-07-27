package http

type PostTrainingUser struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type PostTrainingGroup struct {
	User      PostTrainingUser `json:"user"`
	GroupName string           `json:"group_name"`
	GroupDesc string           `json:"group_desc"`
	Date      string           `json:"date"`
}
