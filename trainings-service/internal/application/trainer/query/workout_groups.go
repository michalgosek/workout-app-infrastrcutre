package query

import (
	"context"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/sirupsen/logrus"
)

type WorkoutGroupsDetails struct {
	WorkoutGroups []WorkoutGroupDetails
}

type WorkoutGroupsHandler struct {
	repository TrainerRepository
}

func (w *WorkoutGroupsHandler) Do(ctx context.Context, trainerUUID string) (WorkoutGroupsDetails, error) {
	logger := logrus.WithFields(logrus.Fields{"Trainer-QRY": "GetWorkoutsHandler"})
	groups, err := w.repository.QueryTrainerWorkoutGroups(ctx, trainerUUID)
	if err != nil {
		logger.Errorf("query workout groups for trainerUUID: %s failed: %v", trainerUUID, err)
		return WorkoutGroupsDetails{}, ErrRepositoryFailure
	}
	var workoutGroups []WorkoutGroupDetails
	for _, g := range groups { // O(n^2)
		workoutGroups = append(workoutGroups, WorkoutGroupDetails{
			TrainerUUID: g.TrainerUUID(),
			TrainerName: g.TrainerName(),
			GroupUUID:   g.UUID(),
			GroupDesc:   g.Description(),
			GroupName:   g.Name(),
			Customers:   ConvertToCustomersData(g.CustomerDetails()),
			Date:        g.Date().String(),
		})
	}
	out := WorkoutGroupsDetails{
		WorkoutGroups: workoutGroups,
	}
	return out, nil
}

func ConvertToCustomersData(details []customer.Details) []CustomerData {
	var customersData []CustomerData
	for _, d := range details {
		customersData = append(customersData, CustomerData{
			UUID: d.UUID(),
			Name: d.Name(),
		})
	}
	return customersData
}

func NewWorkoutGroupsHandler(t TrainerRepository) *WorkoutGroupsHandler {
	if t == nil {
		panic("nil trainer repository")
	}
	return &WorkoutGroupsHandler{
		repository: t,
	}
}
