package testutil

import (
	"sort"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

func GenerateTrainerSchedule(trainerUUID string) trainer.TrainerSchedule {
	ts := time.Now()
	ts = ts.Add(24 * time.Hour)
	name := "dummy"
	desc := "dummy"
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, ts)
	if err != nil {
		panic(err)
	}
	return *schedule
}

func GenerateTrainerSchedules(trainerUUID string, n int) []trainer.TrainerSchedule {
	var schedules []trainer.TrainerSchedule
	for i := 0; i < n; i++ {
		schedule := GenerateTrainerSchedule(trainerUUID)
		schedules = append(schedules, schedule)
	}
	sort.SliceStable(schedules, func(i, j int) bool {
		return schedules[i].UUID() < schedules[j].UUID()
	})
	return schedules
}

func GenerateCustomerSchedule(customerUUID string) customer.CustomerSchedule {
	session, err := customer.NewSchedule(customerUUID)
	if err != nil {
		panic(err)
	}
	return *session
}
