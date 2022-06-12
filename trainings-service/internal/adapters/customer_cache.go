package adapters

import (
	"context"
	"fmt"
	"sync"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
)

type CustomerSchedulesCache struct {
	lookup sync.Map
}

func (c *CustomerSchedulesCache) UpsertSchedule(ctx context.Context, session domain.CustomerWorkoutSession) error {
	c.lookup.Store(session.UserUUID(), session)
	return nil
}

func (c *CustomerSchedulesCache) CancelSchedule(ctx context.Context, sessionUUID, customerUUID string) error {
	customerSessionMapVal, ok := c.lookup.Load(customerUUID)
	if !ok {
		return nil
	}
	customerSession, ok := customerSessionMapVal.(domain.CustomerWorkoutSession)
	if !ok {
		return fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}

	customerSession.UnregisterWorkout(sessionUUID)
	c.lookup.Store(customerSession.UUID(), customerSession)
	return nil
}

func (c *CustomerSchedulesCache) QuerySchedule(_ context.Context, customerUUID string) (domain.CustomerWorkoutSession, error) {
	customerSessionMapVal, ok := c.lookup.Load(customerUUID)
	if !ok {
		return domain.CustomerWorkoutSession{}, nil
	}
	customerSession, ok := customerSessionMapVal.(domain.CustomerWorkoutSession)
	if !ok {
		return domain.CustomerWorkoutSession{}, fmt.Errorf("%w : key: %s", ErrUnderlyingValueType, customerUUID)
	}
	return customerSession, nil
}

func NewCustomerSchedulesCache() *CustomerSchedulesCache {
	return &CustomerSchedulesCache{}
}
