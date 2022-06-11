package domain

import "time"

// Bussiness logic
// User:
// User have 5 workouts to use

// Trainer:
// Trainer cannot have more than 10 people during session and not less than 1
// Training date must be not earlier than 3 hours from current date

type UserWorkoutSession struct {
	UserUUID        string
	Limit           int
	WorkoutSessions []string
}

type TrainerWorkoutSession struct {
	UUID        string
	TrainerUUID string
	Name        string
	Desc        string
	Places      int
	Canceled    bool
	Users       []string
	Date        time.Time
}
