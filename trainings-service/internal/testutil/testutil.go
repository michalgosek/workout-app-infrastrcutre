package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

func getStaticTime() time.Time {
	ts, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	return ts
}

func NewTrainerWorkoutGroup(trainerUUID string) trainer.WorkoutGroup {
	ts := getStaticTime()
	name := "dummy"
	desc := "dummy"
	schedule, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, ts)
	if err != nil {
		panic(err)
	}
	return *schedule
}

func NewWorkoutDay(customerUUID string) customer.WorkoutDay {
	ts := getStaticTime()
	workoutUUID := uuid.NewString()
	session, err := customer.NewWorkoutDay(customerUUID, workoutUUID, ts)
	if err != nil {
		panic(err)
	}
	return *session
}