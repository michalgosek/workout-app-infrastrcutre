package customer

import (
	"errors"
	"fmt"
	"sort"

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

func (c *CustomerSchedule) AssignWorkout(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("%w: into customer workout session", ErrEmptyScheduleUUID)
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
		return ErrSchedulesLimitExceeded
	}

	c.workoutUUIDs = append(c.workoutUUIDs, UUID)
	c.limit--
	return nil
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

func NewSchedule(userUUID string) (*CustomerSchedule, error) {
	// verify logic
	c := CustomerSchedule{
		uuid:     uuid.NewString(),
		userUUID: userUUID,
		limit:    5,
	}
	return &c, nil
}

var ErrEmptyScheduleUUID = errors.New("empty trainer workout UUID")
var ErrSchedulesLimitExceeded = errors.New("customer session workouts number exceeded")
