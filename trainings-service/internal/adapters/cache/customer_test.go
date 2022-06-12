package cache_test

import (
	"context"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/cache"
	"github.com/stretchr/testify/assert"
)

func TestShouldInsertCustomerWorkoutSessionWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)
	expectedSession := customerSession

	SUT := cache.NewCustomerSchedules()

	// when:
	err := SUT.UpsertSchedule(ctx, customerSession)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldUnregisterCustomerWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)
	expectedSession := customerSession

	SUT := cache.NewCustomerSchedules()
	SUT.UpsertSchedule(ctx, customerSession)

	// when:
	err := SUT.CancelSchedule(ctx, customerSession.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldNotReturnErrorAfterUnregisterNonExistingWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)

	SUT := cache.NewCustomerSchedules()

	// when:
	err := SUT.CancelSchedule(ctx, customerSession.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QuerySchedule(ctx, customerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}
