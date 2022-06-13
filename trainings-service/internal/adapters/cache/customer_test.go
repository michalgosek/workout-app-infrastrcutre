package cache

import (
	"context"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldInsertCustomerWorkoutSessionWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSchedule := testutil.GenerateCustomerSchedule(customerUUID)
	expectedSchedule := customerSchedule

	SUT := newCustomerSchedules()

	// when:
	err := SUT.UpsertSchedule(ctx, customerSchedule)

	// then:
	assert.Nil(err)

	actualSchedule, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func TestShouldUnregisterCustomerWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSchedule := testutil.GenerateCustomerSchedule(customerUUID)
	expectedSchedule := customerSchedule

	SUT := newCustomerSchedules()
	SUT.UpsertSchedule(ctx, customerSchedule)

	// when:
	err := SUT.CancelSchedule(ctx, customerSchedule.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSchedule, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSchedule)
}

func TestShouldNotReturnErrorAfterUnregisterNonExistingWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSchedule := testutil.GenerateCustomerSchedule(customerUUID)

	SUT := newCustomerSchedules()

	// when:
	err := SUT.CancelSchedule(ctx, customerSchedule.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSchedule, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Empty(actualSchedule)
}
