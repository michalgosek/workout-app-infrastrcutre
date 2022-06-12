package adapters_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestShouldNotRegisterWorkoutSessionForNonExistingCustomer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	fakeCustomerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	trainerSession := GenerateTestTrainerWorkoutSession("97916cbc-f69b-4602-ac9d-b163b791e73b")
	expectedSession := trainerSession

	SUT := adapters.NewWorkoutsCacheRepoistory()
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
	trainerSession := GenerateTestTrainerWorkoutSession("97916cbc-f69b-4602-ac9d-b163b791e73b")
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)
	expectedSession := customerSession

	SUT := adapters.NewWorkoutsCacheRepoistory()
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
	firstTrainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	secondTrainerSession := GenerateTestTrainerWorkoutSession("ba958504-e56d-438a-8d23-683da191c2f5")

	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)

	expectedSession := customerSession
	expectedSession.AssignWorkout(firstTrainerSession.UUID())
	expectedSession.AssignWorkout(secondTrainerSession.UUID())

	SUT := adapters.NewWorkoutsCacheRepoistory()
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
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)

	expectedSession := customerSession
	expectedSession.AssignWorkout(trainerSession.UUID())

	SUT := adapters.NewWorkoutsCacheRepoistory()
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
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := GenerateTestCustomerWorkoutSession(customerUUID)

	SUT := adapters.NewWorkoutsCacheRepoistory()
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

func GenerateTestTrainerWorkoutSession(trainerUUID string) domain.TrainerWorkoutSession {
	return GenerateTestTrainerWorkoutSessions(trainerUUID, 1)[0]
}

func GenerateTestTrainerWorkoutSessions(trainerUUID string, n int) []domain.TrainerWorkoutSession {
	ts := time.Now()
	ts = ts.Add(24 * time.Hour)

	var sessions []domain.TrainerWorkoutSession
	for i := 0; i < n; i++ {
		name := uuid.NewString()
		desc := uuid.NewString()
		workout, err := domain.NewTrainerWorkoutSession(trainerUUID, name, desc, ts)
		if err != nil {
			panic(err)
		}
		sessions = append(sessions, *workout)
	}

	sort.SliceStable(sessions, func(i, j int) bool {
		return sessions[i].UUID() < sessions[j].UUID()
	})
	return sessions
}

func GenerateTestCustomerWorkoutSession(customerUUID string) domain.CustomerWorkoutSession {
	session, err := domain.NewCustomerWorkoutSessions(customerUUID)
	if err != nil {
		panic(err)
	}
	return *session
}
