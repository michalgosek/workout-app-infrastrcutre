package customer

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutDay struct {
	uuid         string
	customerUUID string
	customerName string
	groupUUID    string
	date         time.Time
}

func (c *WorkoutDay) CustomerName() string {
	return c.customerName
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

func NewWorkoutDay(customerUUID, customerName, groupUUID string, date time.Time) (WorkoutDay, error) {
	if customerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if customerName == "" {
		return WorkoutDay{}, ErrEmptyCustomerName
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
		customerName: customerName,
		date:         date,
	}
	return c, nil
}

func UnmarshalFromDatabase(workoutDayUUID, groupUUID, customerUUID, customerName string, date time.Time) (WorkoutDay, error) {
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
	if customerName == "" {
		return WorkoutDay{}, ErrEmptyCustomerName
	}
	c := WorkoutDay{
		uuid:         workoutDayUUID,
		groupUUID:    groupUUID,
		customerUUID: customerUUID,
		customerName: customerName,
		date:         date,
	}
	return c, nil
}
