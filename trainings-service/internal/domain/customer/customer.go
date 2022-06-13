package customer

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Bussiness logic
// * Customer have 5 workouts to use

type CustomerSchedule struct {
	uuid         string
	userUUID     string
	limit        int
	workoutUUIDs []string
}

func (c *CustomerSchedule) UUID() string {
	return c.uuid
}

func (c *CustomerSchedule) UserUUID() string {
	return c.userUUID
}

func (c *CustomerSchedule) AssignedWorkouts() int {
	return len(c.workoutUUIDs)
}

func (c *CustomerSchedule) WorkoutUUIDs() []string {
	return c.workoutUUIDs
}

func (c *CustomerSchedule) Limit() int {
	return c.limit
}

func (c *CustomerSchedule) UnregisterWorkout(UUID string) {
	var filtered []string
	for _, u := range c.workoutUUIDs {
		if u == UUID {
			continue
		}
		filtered = append(filtered, u)
	}
	c.workoutUUIDs = filtered
}

func (c *CustomerSchedule) AssignWorkout(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyScheduleUUID)
	}
	if c.limit == 0 {
		return ErrSchedulesLimitExceeded
	}
	if len(c.workoutUUIDs) == 0 {
		c.workoutUUIDs = append(c.workoutUUIDs, UUID)
		c.limit--
		return nil
	}
	for _, u := range c.workoutUUIDs {
		if u == UUID {
			return ErrScheduleDuplicate
		}
	}
	c.workoutUUIDs = append(c.workoutUUIDs, UUID)
	c.limit--
	return nil
}

func NewSchedule(userUUID string) (*CustomerSchedule, error) {
	// verify logic
	c := CustomerSchedule{
		uuid:     uuid.NewString(),
		userUUID: userUUID,
		limit:    5,
	}
	return &c, nil
}

var (
	ErrEmptyScheduleUUID      = errors.New("empty trainer workout UUID")
	ErrSchedulesLimitExceeded = errors.New("schedules limit exceeded")
	ErrScheduleDuplicate      = errors.New("schedule duplicate found")
)
