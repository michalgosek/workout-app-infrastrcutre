package domain

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

// Bussiness logic
// * Customer have 5 workouts to use

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

func (c *CustomerWorkoutSession) AssignedWorkouts() int {
	return len(c.workoutUUIDs)
}

func (c *CustomerWorkoutSession) WorkoutUUIDs() []string {
	return c.workoutUUIDs
}

func (c *CustomerWorkoutSession) Limit() int {
	return c.limit
}

func (c *CustomerWorkoutSession) AssignWorkout(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyTrainerWorkoutSessionUUID)
	}
	if len(c.workoutUUIDs) == 0 {
		c.workoutUUIDs = append(c.workoutUUIDs, UUID)
		c.limit--
		return nil
	}

	sort.SliceStable(c.workoutUUIDs, func(i, j int) bool { return c.workoutUUIDs[i] < c.workoutUUIDs[j] })
	idx := sort.Search(len(c.workoutUUIDs)-1, func(i int) bool { return c.workoutUUIDs[i] == UUID })
	if idx < len(c.userUUID) && c.workoutUUIDs[idx] == UUID {
		return nil
	}
	if c.limit == 0 {
		return ErrCustomerWorkouSessionLimitExceeded
	}

	c.workoutUUIDs = append(c.workoutUUIDs, UUID)
	c.limit--
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
		uuid:     uuid.NewString(),
		userUUID: userUUID,
		limit:    5,
	}
	return &c, nil
}

var ErrEmptyTrainerWorkoutSessionUUID = errors.New("empty trainer workout UUID")
var ErrCustomerWorkouSessionLimitExceeded = errors.New("customer session workouts number exceeded")
