package testutil

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	"time"

	"github.com/google/uuid"
)

func getStaticTime() time.Time {
	ts, _ := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	return ts
}

func NewTrainerWorkoutGroup(trainerUUID string) trainings.WorkoutGroup {
	const (
		groupName   = "dummy"
		groupDesc   = "dummy"
		trainerName = "John Doe"
	)
	ts := getStaticTime()
	group, err := trainings.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, ts)
	if err != nil {
		panic(err)
	}
	return group
}

func NewWorkoutDay(customerUUID string) trainings.WorkoutDay {
	ts := getStaticTime()
	name := "John Doe"
	trainerUUID := uuid.NewString()
	groupUUID := uuid.NewString()
	workoutDay, err := trainings.NewWorkoutDay(customerUUID, name, groupUUID, trainerUUID, ts)
	if err != nil {
		panic(err)
	}
	return workoutDay
}
