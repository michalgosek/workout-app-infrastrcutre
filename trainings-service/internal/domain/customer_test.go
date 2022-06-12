package domain_test

import (
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestShouldAssingOneWorkoutWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	customerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	workoutUUID := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	customerWorkoutSession := GenerateTestCustomerWorkoutSession(customerUUID)

	// when:
	err := customerWorkoutSession.AssignWorkout(workoutUUID)

	// then:
	assert.Nil(err)
	assert.Equal(customerWorkoutSession.Limit(), 4)
	assert.Equal(customerWorkoutSession.AssignedWorkouts(), 1)
}

func TestShouldNotAssingDuplicateWorkoutUUID_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	customerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	workoutUUID1 := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	workoutUUID2 := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"

	customerWorkoutSession := GenerateTestCustomerWorkoutSession(customerUUID)

	// when:
	err1 := customerWorkoutSession.AssignWorkout(workoutUUID1)
	err2 := customerWorkoutSession.AssignWorkout(workoutUUID2)

	// then:
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Equal(customerWorkoutSession.Limit(), 4)
}

func TestShouldAssignTwoWorkoutsWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	customerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	workoutUUID1 := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	workoutUUID2 := "cb4bcff9-0e30-4d53-bcd7-87110e786b15"

	customerWorkoutSession := GenerateTestCustomerWorkoutSession(customerUUID)

	// when:
	err1 := customerWorkoutSession.AssignWorkout(workoutUUID1)
	err2 := customerWorkoutSession.AssignWorkout(workoutUUID2)

	// then:
	assert.Nil(err1, err2)
	assert.Equal(customerWorkoutSession.Limit(), 3)
	assert.Equal(customerWorkoutSession.AssignedWorkouts(), 2)
}

func TestShouldNotAssingEmptyWorkoutUUID_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	customerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	customerWorkoutSession := GenerateTestCustomerWorkoutSession(customerUUID)

	// when:
	err := customerWorkoutSession.AssignWorkout("")

	// then:
	assert.ErrorIs(err, domain.ErrEmptyTrainerWorkoutSessionUUID)
}

func TestShouldReturnErrorWhenWorkoutsLimitExeeced_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	customerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	workoutUUID1 := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	workoutUUID2 := "eac9f4e9-8849-4c8c-8235-70429692d2a4"
	workoutUUID3 := "1a4f4c39-e876-4f75-86ff-f93c4862a923"
	workoutUUID4 := "fcd55d25-9734-4357-8a38-ae051219ff2b"
	workoutUUID5 := "7dc57240-9373-47df-9ea7-8876976fcf96"
	workoutUUID6 := "cb4bcff9-0e30-4d53-bcd7-87110e786b15"

	customerWorkoutSession := GenerateTestCustomerWorkoutSession(customerUUID)
	expectedAssginedWorkouts := 5
	expecterdWorkoutLimit := 0

	// when:
	err1 := customerWorkoutSession.AssignWorkout(workoutUUID1)
	err2 := customerWorkoutSession.AssignWorkout(workoutUUID2)
	err3 := customerWorkoutSession.AssignWorkout(workoutUUID3)
	err4 := customerWorkoutSession.AssignWorkout(workoutUUID4)
	err5 := customerWorkoutSession.AssignWorkout(workoutUUID5)
	err6 := customerWorkoutSession.AssignWorkout(workoutUUID6)

	// then:
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Nil(err4)
	assert.Nil(err5)
	assert.ErrorIs(domain.ErrCustomerWorkouSessionLimitExceeded, err6)
	assert.Equal(expecterdWorkoutLimit, customerWorkoutSession.Limit())
	assert.Equal(expectedAssginedWorkouts, customerWorkoutSession.AssignedWorkouts())

}

func GenerateTestCustomerWorkoutSession(customerUUID string) domain.CustomerWorkoutSession {
	c, err := domain.NewCustomerWorkoutSessions(customerUUID)
	if err != nil {
		panic(err)
	}
	return *c
}
