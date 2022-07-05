package query

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

type CustomerData struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type WorkoutGroupDetails struct {
	TrainerUUID string         `json:"trainer_uuid"`
	TrainerName string         `json:"trainer_name"`
	GroupUUID   string         `json:"group_uuid"`
	GroupDesc   string         `json:"group_desc"`
	GroupName   string         `json:"group_name"`
	Customers   []CustomerData `json:"customers"`
	Date        string         `json:"date"`
}

func ConvertToWorkoutGroupsDetails(groups ...trainer.WorkoutGroup) []WorkoutGroupDetails {
	var out []WorkoutGroupDetails
	for _, g := range groups { // O(n^2)
		out = append(out, WorkoutGroupDetails{
			TrainerUUID: g.TrainerUUID(),
			TrainerName: g.TrainerName(),
			GroupUUID:   g.UUID(),
			GroupDesc:   g.Description(),
			GroupName:   g.Name(),
			Customers:   convertToCustomersData(g.CustomerDetails()),
			Date:        g.Date().String(),
		})
	}
	return out
}

func convertToCustomersData(details []customer.Details) []CustomerData {
	var customersData []CustomerData
	for _, d := range details {
		customersData = append(customersData, CustomerData{
			UUID: d.UUID(),
			Name: d.Name(),
		})
	}
	return customersData
}
