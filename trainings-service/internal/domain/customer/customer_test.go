package customer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/stretchr/testify/assert"
)

func TestShouldAssignOneScheduleToCustomerWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const workoutUUID = "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	const scheduleLeft = 4
	const scheduleAssgined = 1

	SUT := GenerateTestcustomerSchedule(customerUUID)

	// when:
	err := SUT.AssignWorkout(workoutUUID)

	// then:
	assert.Nil(err)
	assert.Equal(SUT.Limit(), scheduleLeft)
	assert.Equal(SUT.AssignedWorkouts(), scheduleAssgined)
}

func TestShouldReturnErrorWhenAssignDuplicateScheduleToCustomer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const workoutUUID1 = "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	const workoutUUID2 = "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	const schedulesLeft = 4

	SUT := GenerateTestcustomerSchedule(customerUUID)
	SUT.AssignWorkout(workoutUUID1)

	// when:
	err := SUT.AssignWorkout(workoutUUID2)

	// then:
	assert.Equal(err, customer.ErrScheduleDuplicate)
	assert.Equal(SUT.Limit(), schedulesLeft)
}

func TestShouldAssignTwoSchedulesToCustomerWithSuccess_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const workoutUUID1 = "15939cbe-1f08-4e4a-acf5-47b1bc2e4ad3"
	const workoutUUID2 = "cb4bcff9-0e30-4d53-bcd7-87110e786b15"
	const scheduleAssgined = 2
	const scheduleLeft = 3

	SUT := GenerateTestcustomerSchedule(customerUUID)
	SUT.AssignWorkout(workoutUUID1)

	// when:
	err := SUT.AssignWorkout(workoutUUID2)

	// then:
	assert.Nil(err)
	assert.Equal(SUT.Limit(), scheduleLeft)
	assert.Equal(SUT.AssignedWorkouts(), scheduleAssgined)
}

func TestShouldReturnErrorWhenAssignEmptyScheduleUUIDToCustomer_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	SUT := GenerateTestcustomerSchedule(customerUUID)

	// when:
	err := SUT.AssignWorkout("")

	// then:
	assert.ErrorIs(err, customer.ErrEmptyScheduleUUID)
}

func TestShouldReturnErrorWhenCustomerScheduleLimitExeeced_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const customerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const workoutUUID6 = "cb4bcff9-0e30-4d53-bcd7-87110e786b15"
	const scheduleAssgined = 5
	const scheduleLeft = 0

	SUT := GenerateTestcustomerSchedule(customerUUID)
	AssignScheduleUUIDsToCustomer(&SUT, 5)

	// when:
	err := SUT.AssignWorkout(workoutUUID6)

	// then:
	assert.ErrorIs(customer.ErrSchedulesLimitExceeded, err)
	assert.Equal(scheduleLeft, SUT.Limit())
	assert.Equal(scheduleAssgined, SUT.AssignedWorkouts())
}

func GenerateTestcustomerSchedule(customerUUID string) customer.CustomerSchedule {
	c, err := customer.NewSchedule(customerUUID)
	if err != nil {
		panic(err)
	}
	return *c
}

func AssignScheduleUUIDsToCustomer(schedule *customer.CustomerSchedule, n int) {
	for i := 0; i < n; i++ {
		schedule.AssignWorkout(uuid.NewString())
	}
}
