package domain

import "github.com/google/uuid"

// Bussiness logic
// User:
// User have 5 workouts to use -> the same applies for the lenght of session array

type CustomerWorkoutSession struct {
	uuid         string
	userUUID     string
	limit        int
	workoutUUIDs []string
}

func (c *CustomerWorkoutSession) UUID() string {
	return c.uuid
}

func (c *CustomerWorkoutSession) UserUUID() string {
	return c.userUUID
}

func (c *CustomerWorkoutSession) WorkoutUUIDs() []string {
	return c.workoutUUIDs
}

func (c *CustomerWorkoutSession) Limit() int {
	return c.limit
}

func (c *CustomerWorkoutSession) AssignWorkout(UUIDs ...string) error {
	c.workoutUUIDs = append(c.workoutUUIDs, UUIDs...)
	return nil
}

func (c *CustomerWorkoutSession) UnregisterWorkout(UUID string) {
	var filtered []string
	for _, u := range c.workoutUUIDs {
		if u == UUID {
			continue
		}
		filtered = append(filtered, u)
	}
	c.workoutUUIDs = filtered
}

func NewCustomerWorkoutSessions(userUUID string) (*CustomerWorkoutSession, error) {
	// verify logic
	c := CustomerWorkoutSession{
		uuid:         uuid.NewString(),
		userUUID:     userUUID,
		limit:        5,
		workoutUUIDs: []string{},
	}
	return &c, nil
}
