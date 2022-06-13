package cache_test

import (
	"context"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/cache"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldNotRegisterWorkoutSessionForNonExistingCustomer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	fakeCustomerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	trainerSession := testutil.GenerateTrainerSchedule("97916cbc-f69b-4602-ac9d-b163b791e73b")
	expectedSession := trainerSession

	SUT := cache.NewTrainingSchedules()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)

	// when:
	err := SUT.AssignCustomerToWorkoutSession(ctx, fakeCustomerUUID, trainerSession.UUID())

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestCustomerShouldNotRegisterToNonExistingTrainerWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	fakeSessionUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSession := testutil.GenerateTrainerSchedule("97916cbc-f69b-4602-ac9d-b163b791e73b")
	customerSession := testutil.GenerateCustomerSchedule(customerUUID)
	expectedSession := customerSession

	SUT := cache.NewTrainingSchedules()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)
	SUT.UpsertCustomerWorkoutSession(ctx, customerSession)

	// when:
	err := SUT.AssignCustomerToWorkoutSession(ctx, customerUUID, fakeSessionUUID)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryCustomerWorkoutSession(context.Background(), customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestCustomerShouldBeAssignedToDifferentWorkoutSessionsWithSucccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	firstTrainerSession := testutil.GenerateTrainerSchedule("c27c3952-3bb7-46ce-8700-62906ca192c6")
	secondTrainerSession := testutil.GenerateTrainerSchedule("ba958504-e56d-438a-8d23-683da191c2f5")

	customerSession := testutil.GenerateCustomerSchedule(customerUUID)

	expectedSession := customerSession
	expectedSession.AssignWorkout(firstTrainerSession.UUID())
	expectedSession.AssignWorkout(secondTrainerSession.UUID())

	SUT := cache.NewTrainingSchedules()
	SUT.UpsertTrainerWorkoutSession(ctx, firstTrainerSession)
	SUT.UpsertTrainerWorkoutSession(ctx, secondTrainerSession)
	SUT.UpsertCustomerWorkoutSession(ctx, customerSession)

	// when:
	err1 := SUT.AssignCustomerToWorkoutSession(ctx, customerUUID, firstTrainerSession.UUID())
	err2 := SUT.AssignCustomerToWorkoutSession(ctx, customerUUID, secondTrainerSession.UUID())

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSession, err := SUT.QueryCustomerWorkoutSession(context.Background(), customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestCustomerShouldBeAssignedToTrainerWorkoutSessionWithSucccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	trainerSession := testutil.GenerateTrainerSchedule("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := testutil.GenerateCustomerSchedule(customerUUID)

	expectedSession := customerSession
	expectedSession.AssignWorkout(trainerSession.UUID())

	SUT := cache.NewTrainingSchedules()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)
	SUT.UpsertCustomerWorkoutSession(ctx, customerSession)

	// when:
	err := SUT.AssignCustomerToWorkoutSession(ctx, customerUUID, trainerSession.UUID())

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryCustomerWorkoutSession(context.Background(), customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldRemoveCustomerFromTrainerWorkoutSession_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "b68f3b7e-af79-45d8-ab61-336f5aaff5c8"
	trainerSession := testutil.GenerateTrainerSchedule("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := testutil.GenerateCustomerSchedule(customerUUID)

	SUT := cache.NewTrainingSchedules()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)
	SUT.UpsertCustomerWorkoutSession(ctx, customerSession)
	SUT.AssignCustomerToWorkoutSession(ctx, customerUUID, trainerSession.UUID())

	// when:
	err := SUT.UnregisterCustomerWorkoutSession(ctx, trainerSession.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.NotEmpty(actualSession)
	assert.Equal(actualSession.Customers(), 0)
}
