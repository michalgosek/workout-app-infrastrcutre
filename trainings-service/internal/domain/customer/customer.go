package customer

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutDay struct {
	uuid         string
	customerUUID string
	trainerUUID  string
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

func (c *WorkoutDay) TrainerUUID() string {
	return c.trainerUUID
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

func NewWorkoutDay(customerUUID, customerName, groupUUID string, trainerUUID string, date time.Time) (WorkoutDay, error) {
	if customerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if customerName == "" {
		return WorkoutDay{}, ErrEmptyCustomerName
	}
	if groupUUID == "" {
		return WorkoutDay{}, ErrEmptyGroupUUID
	}
	if trainerUUID == "" {
		return WorkoutDay{}, ErrEmptyTrainerUUID
	}
	if date.IsZero() {
		return WorkoutDay{}, ErrEmptyGroupDate
	}
	c := WorkoutDay{
		uuid:         uuid.NewString(),
		customerUUID: customerUUID,
		trainerUUID:  trainerUUID,
		groupUUID:    groupUUID,
		customerName: customerName,
		date:         date,
	}
	return c, nil
}

type UnmarshalFromDatabaseArgs struct {
	WorkoutDayUUID string
	TrainerUUID    string
	GroupUUID      string
	CustomerUUID   string
	CustomerName   string
	Date           time.Time
}

func UnmarshalFromDatabase(args UnmarshalFromDatabaseArgs) (WorkoutDay, error) {
	if args.WorkoutDayUUID == "" {
		return WorkoutDay{}, ErrEmptyWorkoutDayUUID
	}
	if args.CustomerUUID == "" {
		return WorkoutDay{}, ErrEmptyCustomerUUID
	}
	if args.GroupUUID == "" {
		return WorkoutDay{}, ErrEmptyGroupUUID
	}
	if args.Date.IsZero() {
		return WorkoutDay{}, ErrEmptyGroupDate
	}
	if args.CustomerName == "" {
		return WorkoutDay{}, ErrEmptyCustomerName
	}
	if args.TrainerUUID == "" {
		return WorkoutDay{}, ErrEmptyTrainerUUID
	}
	c := WorkoutDay{
		uuid:         args.WorkoutDayUUID,
		groupUUID:    args.GroupUUID,
		trainerUUID:  args.TrainerUUID,
		customerUUID: args.CustomerUUID,
		customerName: args.CustomerName,
		date:         args.Date,
	}
	return c, nil
}
