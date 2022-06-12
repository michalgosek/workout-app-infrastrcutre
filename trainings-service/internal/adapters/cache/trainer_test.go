package cache_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/cache"

	"github.com/stretchr/testify/assert"
)

func TestTrainerShouldUpdateManyTrainerSchedulesWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	ctx := context.Background()
	SUT := cache.NewTrainerSchedules()

	SUT.UpsertSchedule(ctx, trainerSessions[0])
	SUT.UpsertSchedule(ctx, trainerSessions[1])

	trainerSessions[0].SetName("dummy")
	trainerSessions[1].SetDesc("dummy")
	expectedSchedules := trainerSessions

	// when:
	err1 := SUT.UpsertSchedule(ctx, trainerSessions[0])
	err2 := SUT.UpsertSchedule(ctx, trainerSessions[1])

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSchedules, err := SUT.QuerySchedules(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Equal(expectedSchedules, actualSchedules)
}

func TestTrainerShouldUpdateScheduleWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSession := GenerateTestTrainerWorkoutSession(trainerUUID)
	ctx := context.Background()
	SUT := cache.NewTrainerSchedules()
	SUT.UpsertSchedule(ctx, trainerSession)

	trainerSession.SetName("dummy")
	expectedSchedules := trainerSession

	// when:
	err := SUT.UpsertSchedule(ctx, trainerSession)

	// then:
	assert.Nil(err)

	actualSchedules, err := SUT.QuerySchedule(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSchedules, actualSchedules)
}

func TestQueryNonExistingTrainerScheduleShouldReturnEmptyResult_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	sessionUUID := uuid.NewString()
	SUT := cache.NewTrainerSchedules()

	// when:
	actualSession, err := SUT.QuerySchedule(ctx, sessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(actualSession)
}

func TestShouldInsertTwoSchedulesForTrainerWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	SUT := cache.NewTrainerSchedules()

	// when:
	err1 := SUT.UpsertSchedule(context.Background(), expectedSessions[0])
	err2 := SUT.UpsertSchedule(context.Background(), expectedSessions[1])

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSchedules, err := SUT.QuerySchedules(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.EqualValues(expectedSessions, actualSchedules)
}

func TestShouldUpsertTrainerScheduleForTrainerWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSession := GenerateTestTrainerWorkoutSession(trainerUUID)
	SUT := cache.NewTrainerSchedules()

	// when:
	err := SUT.UpsertSchedule(context.Background(), expectedSession)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QuerySchedule(context.Background(), expectedSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldReturnEmptyScheduleWhenTrainerAttemptsCancelNonExisting_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerSession := GenerateTestTrainerWorkoutSession("c27c3952-3bb7-46ce-8700-62906ca192c6")
	expectedSession := trainerSession
	nonExistingWorkoutSessionUUID := "10b8151c-a686-4fd4-925e-f0a93a41ba50"

	SUT := cache.NewTrainerSchedules()
	SUT.UpsertSchedule(ctx, trainerSession)

	// when:
	deletedSession, err := SUT.CancelSchedule(ctx, nonExistingWorkoutSessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(deletedSession)

	actualSession, err := SUT.QuerySchedule(context.Background(), trainerSession.UUID())
	assert.Nil(err)
	assert.Equal(expectedSession, actualSession)
}

func TestShouldCancelAllTrainerSchedulesWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSessions := GenerateTestTrainerWorkoutSessions(trainerUUID, 2)
	expectedSessions := trainerSessions
	workoutSessionsUUIDs := []string{trainerSessions[0].UUID(), trainerSessions[1].UUID()}
	SUT := cache.NewTrainerSchedules()

	SUT.UpsertSchedule(ctx, trainerSessions[0])
	SUT.UpsertSchedule(ctx, trainerSessions[1])

	// when:
	deletedSession, err := SUT.CancelSchedules(ctx, workoutSessionsUUIDs...)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSessions, deletedSession)

	actualSession, err := SUT.QuerySchedules(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}
