package customer

import (
	"errors"

	"github.com/google/uuid"
)

// Bussiness logic
// * Customer have 5 workouts to use

type CustomerSchedule struct {
	uuid          string
	customerUUID  string
	limit         int
	scheduleUUIDs []string
}

func (c *CustomerSchedule) UUID() string {
	return c.uuid
}

func (c *CustomerSchedule) CustomerUUID() string {
	return c.customerUUID
}

func (c *CustomerSchedule) AssignedSchedules() int {
	return len(c.scheduleUUIDs)
}

func (c *CustomerSchedule) ScheduleUUIDs() []string {
	return c.scheduleUUIDs
}

func (c *CustomerSchedule) Limit() int {
	return c.limit
}

func (c *CustomerSchedule) UnregisterSchedule(UUID string) {
	var filtered []string
	for _, u := range c.scheduleUUIDs {
		if u == UUID {
			continue
		}
		filtered = append(filtered, u)
	}
	c.scheduleUUIDs = filtered
}

func (c *CustomerSchedule) AssignSchedule(UUID string) error {
	if UUID == "" {
		return ErrEmptyScheduleUUID
	}
	if c.limit == 0 {
		return ErrSchedulesLimitExceeded
	}
	if len(c.scheduleUUIDs) == 0 {
		c.scheduleUUIDs = append(c.scheduleUUIDs, UUID)
		c.limit--
		return nil
	}
	for _, u := range c.scheduleUUIDs {
		if u == UUID {
			return ErrScheduleDuplicate
		}
	}
	c.scheduleUUIDs = append(c.scheduleUUIDs, UUID)
	c.limit--
	return nil
}

func NewSchedule(UUID string) (*CustomerSchedule, error) {
	if UUID == "" {
		return nil, ErrEmptyCustomerUUID
	}
	c := CustomerSchedule{
		uuid:         uuid.NewString(),
		customerUUID: UUID,
		limit:        5,
	}
	return &c, nil
}

func UnmarshalFromDatabase(scheduleUUID, customerUUID string, limit int, scheduleUUIDs []string) CustomerSchedule {
	return CustomerSchedule{
		uuid:          scheduleUUID,
		customerUUID:  customerUUID,
		limit:         limit,
		scheduleUUIDs: scheduleUUIDs,
	}
}

var (
	ErrEmptyCustomerUUID      = errors.New("empty customer UUID")
	ErrEmptyScheduleUUID      = errors.New("empty schedule UUID")
	ErrSchedulesLimitExceeded = errors.New("schedules limit exceeded")
	ErrScheduleDuplicate      = errors.New("schedule duplicate found")
)
