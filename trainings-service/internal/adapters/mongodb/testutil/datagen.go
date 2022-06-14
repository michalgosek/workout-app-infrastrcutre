package testutil

import (
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

func GenerateTrainerSchedule(trainerUUID string) trainer.TrainerSchedule {
	ts, err := time.Parse("2006-01-02 15:04", "2099-12-12 23:30")
	if err != nil {
		panic(err)
	}
	name := "dummy"
	desc := "dummy"
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, ts)
	if err != nil {
		panic(err)
	}
	return *schedule
}

func GenerateCustomerSchedule(customerUUID string) customer.CustomerSchedule {
	session, err := customer.NewSchedule(customerUUID)
	if err != nil {
		panic(err)
	}
	return *session
}
