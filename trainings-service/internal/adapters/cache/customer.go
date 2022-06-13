package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

type CustomerSchedules struct {
	lookup sync.Map
}

func (c *CustomerSchedules) UpsertSchedule(ctx context.Context, schedule domain.CustomerWorkoutSession) error {
	c.lookup.Store(schedule.UserUUID(), schedule)
	return nil
}

func (c *CustomerSchedules) CancelSchedule(ctx context.Context, scheduleUUID, customerUUID string) error {
	customerScheduleMapV, ok := c.lookup.Load(customerUUID)
	if !ok {
		return nil
	}
	schedule, ok := customerScheduleMapV.(domain.CustomerWorkoutSession)
	if !ok {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}
	schedule.UnregisterWorkout(scheduleUUID)
	c.lookup.Store(schedule.UUID(), schedule)
	return nil
}

func (c *CustomerSchedules) QuerySchedule(_ context.Context, customerUUID string) (domain.CustomerWorkoutSession, error) {
	customerScheduleMapV, ok := c.lookup.Load(customerUUID)
	if !ok {
		return domain.CustomerWorkoutSession{}, nil
	}
	schedule, ok := customerScheduleMapV.(domain.CustomerWorkoutSession)
	if !ok {
		return domain.CustomerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}
	return schedule, nil
}

func newCustomerSchedules() *CustomerSchedules {
	return &CustomerSchedules{}
}
