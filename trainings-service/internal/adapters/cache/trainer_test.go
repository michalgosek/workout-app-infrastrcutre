package cache

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"

	"github.com/stretchr/testify/assert"
)

func TestTrainerShouldUpdateManyTrainerSchedulesWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSchedules := testutil.GenerateTrainerSchedules(trainerUUID, 2)
	ctx := context.Background()
	SUT := newTrainerSchedules()

	SUT.UpsertSchedule(ctx, trainerSchedules[0])
	SUT.UpsertSchedule(ctx, trainerSchedules[1])

	trainerSchedules[0].SetName("dummy")
	trainerSchedules[1].SetDesc("dummy")
	expectedSchedules := trainerSchedules

	// when:
	err1 := SUT.UpsertSchedule(ctx, trainerSchedules[0])
	err2 := SUT.UpsertSchedule(ctx, trainerSchedules[1])

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
	trainerSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	ctx := context.Background()
	SUT := newTrainerSchedules()
	SUT.UpsertSchedule(ctx, trainerSchedule)

	trainerSchedule.SetName("dummy")
	expectedSchedules := trainerSchedule

	// when:
	err := SUT.UpsertSchedule(ctx, trainerSchedule)

	// then:
	assert.Nil(err)

	actualSchedules, err := SUT.QuerySchedule(context.Background(), trainerSchedule.UUID())
	assert.Nil(err)
	assert.Equal(expectedSchedules, actualSchedules)
}

func TestQueryNonExistingTrainerScheduleShouldReturnEmptyResult_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	sessionUUID := uuid.NewString()
	SUT := newTrainerSchedules()

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
	expectedSchedules := testutil.GenerateTrainerSchedules(trainerUUID, 2)
	SUT := newTrainerSchedules()

	// when:
	err1 := SUT.UpsertSchedule(context.Background(), expectedSchedules[0])
	err2 := SUT.UpsertSchedule(context.Background(), expectedSchedules[1])

	// then:
	assert.Nil(err1)
	assert.Nil(err2)

	actualSchedules, err := SUT.QuerySchedules(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.EqualValues(expectedSchedules, actualSchedules)
}

func TestShouldUpsertTrainerScheduleForTrainerWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	expectedSchedule := testutil.GenerateTrainerSchedule(trainerUUID)
	SUT := newTrainerSchedules()

	// when:
	err := SUT.UpsertSchedule(context.Background(), expectedSchedule)

	// then:
	assert.Nil(err)

	actualSession, err := SUT.QuerySchedule(context.Background(), expectedSchedule.UUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSession)
}

func TestShouldReturnEmptyScheduleWhenTrainerAttemptsCancelNonExisting_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerSchedule := testutil.GenerateTrainerSchedule("c27c3952-3bb7-46ce-8700-62906ca192c6")
	expectedSchedule := trainerSchedule
	nonExistingWorkoutSessionUUID := "10b8151c-a686-4fd4-925e-f0a93a41ba50"

	SUT := newTrainerSchedules()
	SUT.UpsertSchedule(ctx, trainerSchedule)

	// when:
	deletedSession, err := SUT.CancelSchedule(ctx, nonExistingWorkoutSessionUUID)

	// then:
	assert.Nil(err)
	assert.Empty(deletedSession)

	actualSession, err := SUT.QuerySchedule(context.Background(), trainerSchedule.UUID())
	assert.Nil(err)
	assert.Equal(expectedSchedule, actualSession)
}

func TestShouldCancelAllTrainerSchedulesWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	ctx := context.Background()
	trainerUUID := "c27c3952-3bb7-46ce-8700-62906ca192c6"
	trainerSchedules := testutil.GenerateTrainerSchedules(trainerUUID, 2)
	expectedSchedules := trainerSchedules
	scheduleUUIDs := []string{trainerSchedules[0].UUID(), trainerSchedules[1].UUID()}
	SUT := newTrainerSchedules()

	SUT.UpsertSchedule(ctx, trainerSchedules[0])
	SUT.UpsertSchedule(ctx, trainerSchedules[1])

	// when:
	deletedSession, err := SUT.CancelSchedules(ctx, scheduleUUIDs...)

	// then:
	assert.Nil(err)
	assert.Equal(expectedSchedules, deletedSession)

	actualSession, err := SUT.QuerySchedules(context.Background(), trainerUUID)
	assert.Nil(err)
	assert.Empty(actualSession)
}
