package customer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WorkoutDay struct {
	uuid                    string
	customerUUID            string
	trainerWorkoutGroupUUID string
	date                    time.Time
}

func (c *WorkoutDay) Date() time.Time {
	return c.date
}

func (c *WorkoutDay) TrainerWorkoutGroupUUID() string {
	return c.trainerWorkoutGroupUUID
}

func (c *WorkoutDay) UUID() string {
	return c.uuid
}

func (c *WorkoutDay) CustomerUUID() string {
	return c.customerUUID
}

func NewWorkoutDay(customerUUID, trainerWorkoutGroupUUID string, date time.Time) (*WorkoutDay, error) {
	if customerUUID == "" {
		return nil, ErrEmptyCustomerUUID
	}
	if trainerWorkoutGroupUUID == "" {
		return nil, ErrEmptyTrainerWorkoutGroupUUID
	}
	if date.IsZero() {
		return nil, ErrEmptyWorkoutDate
	}
	c := WorkoutDay{
		uuid:                    uuid.NewString(),
		customerUUID:            customerUUID,
		trainerWorkoutGroupUUID: trainerWorkoutGroupUUID,
		date:                    date,
	}
	return &c, nil
}

func UnmarshalFromDatabase(customerWorkoutDayUUID, trainerWorkoutGroupUUID, customerUUID string, date time.Time) (WorkoutDay, error) {
	if customerWorkoutDayUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerWorkoutUUID
	}
	if customerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if trainerWorkoutGroupUUID == "" {
		return WorkoutDay{}, ErrEmptyTrainerWorkoutGroupUUID
	}
	if date.IsZero() {
		return WorkoutDay{}, ErrEmptyWorkoutDate
	}
	c := WorkoutDay{
		uuid:                    customerWorkoutDayUUID,
		trainerWorkoutGroupUUID: trainerWorkoutGroupUUID,
		customerUUID:            customerUUID,
		date:                    date,
	}
	return c, nil
}

var (
	ErrEmptyCustomerWorkoutUUID     = errors.New("empty customer workout day UUID")
	ErrEmptyWorkoutDate             = errors.New("empty workout date")
	ErrEmptyCustomerUUID            = errors.New("empty customer UUID")
	ErrEmptyTrainerWorkoutGroupUUID = errors.New("empty trainer workout group UUID")
)
