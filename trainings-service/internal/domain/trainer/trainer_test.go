package trainer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

// Trainer:
// Trainer cannot have more than 10 people during session and not less than 1 (same applies to places) X
// Training date must be not earlier than 3 hours from current date X

// Desc cannot be length than 100 chars X
// Name cannot be length than 15 chars X
// places cannot be less than 0 or gerater than userUUIDS
// Add tests for unregister customer

func TestShouldReturnErrorWhenCustomerLimitExeeced_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const customerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const customersLeft = 0
	const customersAssigned = 10

	SUT := testutil.GenerateTrainerSchedule(trainerUUID)
	AssignCustomerToTrainerSchedule(&SUT, 10)

	// when:
	err := SUT.AssignCustomer(customerUUID)

	// then:
	assert.ErrorIs(trainer.ErrCustomersScheduleLimitExceeded, err)
	assert.Equal(customersLeft, SUT.Limit())
	assert.Equal(customersAssigned, SUT.Customers())
}

func AssignCustomerToTrainerSchedule(schedule *trainer.TrainerSchedule, n int) {
	for i := 0; i < n; i++ {
		schedule.AssignCustomer(uuid.NewString())
	}
}
