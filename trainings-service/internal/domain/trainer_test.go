package domain_test

import (
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Trainer:
// Trainer cannot have more than 10 people during session and not less than 1 (same applies to places) X
// Training date must be not earlier than 3 hours from current date X

// Desc cannot be length than 100 chars
// Name cannot be length than 15 chars
// places cannot be less than 0 or gerater than userUUIDS
// Add tests for unregister customer

func TestShouldReturnErrorWhenCustomerLimitExeeced_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	trainerUUID := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	SUT := GenerateTestTrainerWorkoutSession(trainerUUID)

	customerUUID1 := "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	customerUUID2 := "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	customerUUID3 := "eac9f4e9-8849-4c8c-8235-70429692d2a4"
	customerUUID4 := "1a4f4c39-e876-4f75-86ff-f93c4862a923"
	customerUUID5 := "fcd55d25-9734-4357-8a38-ae051219ff2b"
	customerUUID6 := "7dc57240-9373-47df-9ea7-8876976fcf96"
	customerUUID7 := "cb4bcff9-0e30-4d53-bcd7-87110e786b15"
	customerUUID8 := "b9619d7a-9cc4-4ab4-a5fd-745678f160fc"
	customerUUID9 := "ea23e9c9-3720-4a7f-a9c9-5a41324eaeee"
	customerUUID10 := "4e2d91cb-c2c8-4178-91d4-44e70c57edac"
	customerUUID11 := "44403123-40d0-4085-8b6a-8f804e90f044"

	expecterdWorkoutLimit := 0
	expectedAssginedWorkouts := 10

	// when:
	err1 := SUT.AssignCustomer(customerUUID1)
	err2 := SUT.AssignCustomer(customerUUID2)
	err3 := SUT.AssignCustomer(customerUUID3)
	err4 := SUT.AssignCustomer(customerUUID4)
	err5 := SUT.AssignCustomer(customerUUID5)
	err6 := SUT.AssignCustomer(customerUUID6)
	err7 := SUT.AssignCustomer(customerUUID7)
	err8 := SUT.AssignCustomer(customerUUID8)
	err9 := SUT.AssignCustomer(customerUUID9)
	err10 := SUT.AssignCustomer(customerUUID10)
	err11 := SUT.AssignCustomer(customerUUID11)

	// then:
	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.Nil(err4)
	assert.Nil(err5)
	assert.Nil(err6)
	assert.Nil(err7)
	assert.Nil(err8)
	assert.Nil(err9)
	assert.Nil(err10)
	assert.ErrorIs(domain.ErrCustomerWorkouSessionLimitExceeded, err11)
	assert.Equal(expecterdWorkoutLimit, SUT.Limit())
	assert.Equal(expectedAssginedWorkouts, SUT.AssignedCustomers())
}

func GenerateTestTrainerWorkoutSession(trainerUUID string) domain.TrainerWorkoutSession {
	return GenerateTestTrainerWorkoutSessions(trainerUUID, 1)[0]
}

func GenerateTestTrainerWorkoutSessions(trainerUUID string, n int) []domain.TrainerWorkoutSession {
	var sessions []domain.TrainerWorkoutSession
	for i := 0; i < n; i++ {
		name := uuid.NewString()
		desc := uuid.NewString()
		ts := time.Now().Add(24 * time.Hour)
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
