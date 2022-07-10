package query

import (
	"context"
	"errors"
	"fmt"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	"time"
)

type CustomerWorkoutGroupDetails struct {
	TrainerUUID  string
	TrainerName  string
	WorkoutName  string
	WorkoutDesc  string
	Date         time.Time
	Participants int
}

type WorkoutDayHandlerArgs struct {
	GroupUUID    string
	CustomerUUID string
	TrainerUUID  string
}

type WorkoutDayHandler struct {
	trainerService TrainerService
}

func (w *WorkoutDayHandler) Do(ctx context.Context, args WorkoutDayHandlerArgs) (CustomerWorkoutGroupDetails, error) {
	workoutGroup, err := w.trainerService.GetCustomerWorkoutGroup(ctx, trainer.WorkoutGroupWithCustomerArgs{
		GroupUUID:    args.GroupUUID,
		TrainerUUID:  args.TrainerUUID,
		CustomerUUID: args.CustomerUUID,
	})
	if err != nil {
		return CustomerWorkoutGroupDetails{}, fmt.Errorf("trainer service failure: %w", err)
	}
	customerWorkoutDay := CustomerWorkoutGroupDetails{
		Date:         workoutGroup.Date(),
		TrainerUUID:  workoutGroup.TrainerUUID(),
		TrainerName:  workoutGroup.TrainerName(),
		WorkoutName:  workoutGroup.Name(),
		WorkoutDesc:  workoutGroup.Description(),
		Participants: workoutGroup.AssignedCustomers(),
	}
	return customerWorkoutDay, nil
}

func NewWorkoutDayHandler(t TrainerService) (*WorkoutDayHandler, error) {
	if t == nil {
		return nil, errors.New("nil trainer service")
	}
	h := WorkoutDayHandler{trainerService: t}
	return &h, nil
}
