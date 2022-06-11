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

// Uzytkownik powinien sie zapisac na trening X

// Uzytkownik nie powinien sie zapisac na nie istniejacy trening

// Uzytkownik powinien sie wypisac z treiningu
// Powinno zwrocic nil error, gdy wypisujemy sie z nie istniejacego treningu

func TestShouldInsertCustomerWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "eabfc05c-5a32-42bf-b942-65885a673151"
	customerSession := CreatCustomerWorkoutSessions(customerUUID)
	expectedSession := customerSession

	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	err := SUT.UpsertCustomerWorkoutSession(ctx, customerSession)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryCustomerWorkoutSession(ctx, customerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestCustomerShouldBeAssignedToTrainerWorkoutSessionWithSucccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "e5bf08ee-e287-40b6-8ef7-d43ec18fd17a"
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := CreatCustomerWorkoutSessions(customerUUID)

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

func TestDeleteWorkoutSessionShouldReturnEmptyResultWhenSessionNotExistUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	expectedSession := trainerSession
	nonExistingWorkoutSessionUUID := "10b8151c-a686-4fd4-925e-f0a93a41ba50"

	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)

	// when:
	deletedSession, err := SUT.DeleteTrainerWorkoutSession(ctx, nonExistingWorkoutSessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(deletedSession)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldRemoveCustomerFromTrainerWorkoutSessionUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	customerUUID := "b68f3b7e-af79-45d8-ab61-336f5aaff5c8"
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	customerSession := CreatCustomerWorkoutSessions(customerUUID)

	trainerSession.AssignCustomer(customerUUID)
	customerSession.AssignWorkout(trainerSession.UUID())

	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)
	SUT.UpsertCustomerWorkoutSession(ctx, customerSession)

	// when:
	err := SUT.RemoveCustomerFromTrainerWorkoutSession(ctx, trainerSession.UUID(), customerUUID)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.NotEmpty(actualSession)
	assert.Equal(actualSession.Customers(), 0)
}

func TestTrainerShouldDeleteWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	expectedSession := trainerSession
	workoutSessionUUID := trainerSession.UUID()
	trainerUUID := trainerSession.TrainerUUID()
	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)

	// when:
	deletedSession, err := SUT.DeleteTrainerWorkoutSession(ctx, workoutSessionUUID)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSession, deletedSession)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}

func TestTrainerShouldDeleteWorkoutSessionsWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	expectedSessions := trainerSessions
	workoutSessionsUUIDs := []string{trainerSessions[0].UUID(), trainerSessions[1].UUID()}
	SUT := adapters.NewWorkoutsCacheRepoistory()

	SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[0])
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[1])

	// when:
	deletedSession, err := SUT.DeleteTrainerWorkoutSessions(ctx, workoutSessionsUUIDs...)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSessions, deletedSession)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}

func TestTrainerShouldInsertWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSession := GenerateTestTrainerWorkoutSession(trainerUUID)
	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	err := SUT.UpsertTrainerWorkoutSession(context.Background(), expectedSession)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), expectedSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestTrainerShouldInsertTwoWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	err1 := SUT.UpsertTrainerWorkoutSession(context.Background(), expectedSessions[0])
	err2 := SUT.UpsertTrainerWorkoutSession(context.Background(), expectedSessions[1])

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSessions, err := SUT.QueryTrainerWorkoutSessions(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.EqualValues(expectedSessions, actualSessions)
}

func TestQueryWorkoutSessionShouldReturnEmptyResultWhenSessionNotUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	sessionUUID := uuid.NewString()
	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	actualSession, err := SUT.QueryTrainerWorkoutSession(ctx, sessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSession)
}

func TestTrainerShouldUpdateWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSession := GenerateTestTrainerWorkoutSession(trainerUUID)
	ctx := context.Background()
	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)

	trainerSession.SetName("dummy")
	expectedWorkoutSessions := trainerSession

	// when:
	err := SUT.UpsertTrainerWorkoutSession(ctx, trainerSession)

	// then:
	assert.Nil(err)

	actualSessions, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedWorkoutSessions, actualSessions)
}

func TestTrainerShouldUpdateWorkoutSessionsWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	ctx := context.Background()
	SUT := adapters.NewWorkoutsCacheRepoistory()

	SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[0])
	SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[1])

	trainerSessions[0].SetName("dummy")
	trainerSessions[1].SetDesc("dummy")
	expectedWorkoutSessions := trainerSessions

	// when:
	err1 := SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[0])
	err2 := SUT.UpsertTrainerWorkoutSession(ctx, trainerSessions[1])

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSessions, err := SUT.QueryTrainerWorkoutSessions(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Equal(expectedWorkoutSessions, actualSessions)
}

func GenerateTestTrainerWorkoutSession(trainerUUID string) domain.TrainerWorkoutSession {
	return GenerateTestTrainerWorkoutSessions(trainerUUID, 1)[0]
}

func GenerateTestTrainerWorkoutSessions(trainerUUID string, n int) []domain.TrainerWorkoutSession {
	ts := time.Now()
	ts.Add(3 * time.Hour)

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

func CreatCustomerWorkoutSessions(customerUUID string) domain.CustomerWorkoutSession {
	session, err := domain.NewCustomerWorkoutSessions(customerUUID)
	if err != nil {
		panic(err)
	}
	return *session
}
