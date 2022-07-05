package customer

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutDay struct {
	uuid         string
	customerUUID string
	groupUUID    string
	date         time.Time
}

func (c *WorkoutDay) Date() time.Time {
	return c.date
}

func (c *WorkoutDay) GroupUUID() string {
	return c.groupUUID
}

func (c *WorkoutDay) UUID() string {
	return c.uuid
}

func (c *WorkoutDay) CustomerUUID() string {
	return c.customerUUID
}

func NewWorkoutDay(customerUUID, groupUUID string, date time.Time) (WorkoutDay, error) {
	if customerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if groupUUID == "" {
		return WorkoutDay{}, ErrEmptyGroupUUID
	}
	if date.IsZero() {
		return WorkoutDay{}, ErrEmptyGroupDate
	}
	c := WorkoutDay{
		uuid:         uuid.NewString(),
		customerUUID: customerUUID,
		groupUUID:    groupUUID,
		date:         date,
	}
	return c, nil
}

func UnmarshalFromDatabase(workoutDayUUID, groupUUID, customerUUID string, date time.Time) (WorkoutDay, error) {
	if workoutDayUUID == "" {
		return WorkoutDay{}, ErrEmptyWorkoutDayUUID
	}
	if customerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if groupUUID == "" {
		return WorkoutDay{}, ErrEmptyGroupUUID
	}
	if date.IsZero() {
		return WorkoutDay{}, ErrEmptyGroupDate
	}
	c := WorkoutDay{
		uuid:         workoutDayUUID,
		groupUUID:    groupUUID,
		customerUUID: customerUUID,
		date:         date,
	}
	return c, nil
}
