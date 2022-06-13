package testutil

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

func GenerateTrainerSchedule(trainerUUID string) domain.TrainerSchedule {
	ts := time.Now()
	ts = ts.Add(24 * time.Hour)
	name := uuid.NewString()
	desc := uuid.NewString()
	schedule, err := domain.NewTrainerSchedule(trainerUUID, name, desc, ts)
	if err != nil {
		panic(err)
	}
	return *schedule
}

func GenerateTrainerSchedules(trainerUUID string, n int) []domain.TrainerSchedule {
	var schedules []domain.TrainerSchedule
	for i := 0; i < n; i++ {
		schedule := GenerateTrainerSchedule(trainerUUID)
		schedules = append(schedules, schedule)
	}
	sort.SliceStable(schedules, func(i, j int) bool {
		return schedules[i].UUID() < schedules[j].UUID()
	})
	return schedules
}

func GenerateCustomerSchedule(customerUUID string) domain.CustomerWorkoutSession {
	session, err := domain.NewCustomerWorkoutSessions(customerUUID)
	if err != nil {
		panic(err)
	}
	return *session
}
