package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Trener moze dodac sesje
// Trener moze dodac dwie sesje
// Powinno zwrocic pusta odpowiedz i bez bledu gdy sesji nie ma
// Trener moze zaktualizowac sesje
// Trener moze zaktualizoawc dwie sesje
// Trener moze usunac sesje
// Trener moze usunac wiele sesji
// Powinno zwrocic pusta odpowiedz i bez bledu gdy sesji nie ma w trakcie usuwania

func TestDeleteWorkoutSessionShouldReturnEmptyResultWhenSessionNotExistUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	workouts := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6", 1)
	expectedSession := workouts

	nonExistingWorkoutSessionUUID := "10b8151c-a686-4fd4-925e-f0a93a41ba50"
	trainerUUID := workouts[0].TrainerUUID

	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSessions(ctx, workouts...)

	// when:
	deletedSession, err := SUT.DeleteTrainerWorkoutSession(ctx, nonExistingWorkoutSessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(deletedSession)

	actualSession, err := SUT.QueryTrainerWorkoutSessions(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestTrainerShouldDeleteWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	workouts := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6", 1)

	expectedSession := workouts[0]
	workoutSessionUUID := workouts[0].UUID
	trainerUUID := workouts[0].TrainerUUID

	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSessions(ctx, workouts...)

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
	workouts := GenerateTestTrainerWorkoutSession(trainerUUID, 2)

	expectedSession := workouts
	workoutSessionsUUIDs := []string{workouts[0].UUID, workouts[1].UUID}

	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSessions(ctx, workouts...)

	// when:
	deletedSession, err := SUT.DeleteTrainerWorkoutSessions(ctx, workoutSessionsUUIDs...)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSession, deletedSession)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}

func TestTrainerShouldInsertWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSession := GenerateTestTrainerWorkoutSession(trainerUUID, 1)

	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	err := SUT.UpsertTrainerWorkoutSessions(context.Background(), expectedSession...)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QueryTrainerWorkoutSession(context.Background(), expectedSession[0].UUID)
	assert.Nil(err)
	assert.Equal(expectedSession[0], actualSession)
}

func TestTrainerShouldInsertTwoWorkoutSessionWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSessions := GenerateTestTrainerWorkoutSession(trainerUUID, 2)

	SUT := adapters.NewWorkoutsCacheRepoistory()

	// when:
	err := SUT.UpsertTrainerWorkoutSessions(context.Background(), expectedSessions...)

	// then:
	assert.Nil(err)

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
	workoutSessions := GenerateTestTrainerWorkoutSession(trainerUUID, 1)
	ctx := context.Background()
	SUT := adapters.NewWorkoutsCacheRepoistory()

	SUT.UpsertTrainerWorkoutSessions(ctx, workoutSessions...)

	workoutSessions[0].Canceled = true
	expectedWorkoutSessions := workoutSessions

	// when:
	err := SUT.UpsertTrainerWorkoutSessions(ctx, workoutSessions...)

	// then:
	assert.Nil(err)

	actualSessions, err := SUT.QueryTrainerWorkoutSessions(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Equal(expectedWorkoutSessions, actualSessions)
}

func TestTrainerShouldUpdateWorkoutSessionsWithSuccessUnit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	workoutSession := GenerateTestTrainerWorkoutSession(trainerUUID, 2)
	ctx := context.Background()
	SUT := adapters.NewWorkoutsCacheRepoistory()
	SUT.UpsertTrainerWorkoutSessions(ctx, workoutSession...)

	workoutSession[0].Canceled = true
	workoutSession[1].Name = "dummy"
	expectedWorkoutSessions := workoutSession

	// when:
	err := SUT.UpsertTrainerWorkoutSessions(ctx, expectedWorkoutSessions...)

	// then:
	assert.Nil(err)

	actualSessions, err := SUT.QueryTrainerWorkoutSessions(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Equal(expectedWorkoutSessions, actualSessions)
}

func GenerateTestTrainerWorkoutSession(trainerUUID string, n int) []domain.TrainerWorkoutSession {
	var sessions []domain.TrainerWorkoutSession
	for i := 0; i < n; i++ {
		sessions = append(sessions, domain.TrainerWorkoutSession{
			UUID:        uuid.NewString(),
			TrainerUUID: trainerUUID,
			Name:        uuid.NewString(),
			Desc:        uuid.NewString(),
			Places:      10,
			Canceled:    false,
			Users:       []string{uuid.NewString(), uuid.NewString()},
			Date:        time.Time{},
		})
	}
	return sessions
}
