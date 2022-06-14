package testutil

import (
	"github.com/google/uuid"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

func GenerateTrainerWorkoutGroup(trainerUUID string) trainer.WorkoutGroup {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	name := "dummy"
	desc := "dummy"
	schedule, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, ts)
	if err != nil {
		panic(err)
	}
	return *schedule
}

func GenerateNewWorkoutDay(customerUUID string) customer.WorkoutDay {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	workoutUUID := uuid.NewString()
	session, err := customer.NewWorkoutDay(customerUUID, workoutUUID, ts)
	if err != nil {
		panic(err)
	}
	return *session
}
