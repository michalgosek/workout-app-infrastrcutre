package trainer_test

import (
	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestShouldReturnErrorWhenCustomerLimitExeeced_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const customerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const customersLeft = 0
	const customersAssigned = 10

	SUT := testutil.GenerateTrainerWorkoutGroup(trainerUUID)
	AssignCustomerToTrainerSchedule(&SUT, 10)

	// when:
	err := SUT.AssignCustomer(customerUUID)

	// then:
	assertions.ErrorIs(trainer.ErrCustomersScheduleLimitExceeded, err)
	assertions.Equal(customersLeft, SUT.Limit())
	assertions.Equal(customersAssigned, SUT.AssignedCustomers())
}

func AssignCustomerToTrainerSchedule(workoutGroup *trainer.WorkoutGroup, n int) {
	for i := 0; i < n; i++ {
		workoutGroup.AssignCustomer(uuid.NewString())
	}
}

func TestShouldNotReturnErrorWhenScheduleNameIsUnderLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	const newName = "dummy1"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assertions.Nil(err)
	assertions.Equal(newName, SUT.Name())

}

func TestShouldReturnErrorWhenScheduleNameIsOverLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	newName := strings.Repeat("s", 16)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assertions.Equal(trainer.ErrScheduleNameExceeded, err)
}

func TestShouldNotReturnErrorWhenScheduleNameIsEqualLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	newName := strings.Repeat("s", 15)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assertions.Nil(err)
}

func TestShouldRegisterCustomerToScheduleWithSucces(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 9
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.AssignCustomer(customerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(customersLeft, SUT.Limit())
}

func TestShouldUnregisterCustomerFromScheduleWithSucces(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 10
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	SUT.AssignCustomer(customerUUID)

	// when:
	SUT.UnregisterCustomer(customerUUID)

	// then:
	assertions.Empty(SUT.AssignedCustomers())
	assertions.Equal(customersLeft, SUT.Limit())
}

func TestShouldNotReturnErrorWhenTextLengthIsUnderLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	const newDesc = "dummy1"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(newDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(newDesc, SUT.Desc())
}

func TestShouldReturnErrorWhenTextLengthIsOverLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	invalidDesc := strings.Repeat("a", 101)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(invalidDesc)

	// then:
	assert.ErrorIs(err, trainer.ErrScheduleDescriptionExceeded)
	assert.Equal(desc, SUT.Desc())
}

func TestShouldReturnErrorWhenTextLengthEqualsLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	expectedDesc := strings.Repeat("a", 100)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(expectedDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedDesc, SUT.Desc())
}

func TestShouldReturnErrorWhenSpecifiedTimeIsOneMinEarlierFromNow_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"

	threshold := 3 * time.Hour
	hourEarlier := time.Now().Add(threshold - 1*time.Hour)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, hourEarlier)

	// then:
	assertions.Equal(trainer.ErrScheduleDateViolation, err)
	assertions.Nil(workoutGroup)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneMinLaterFromThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Hour)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, minLater)

	// then:
	assertions.NotNil(workoutGroup)
	assertions.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneSecondLaterFromThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Second)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, minLater)

	// then:
	assertions.NotNil(workoutGroup)
	assertions.Nil(err)
}

func TestShouldReturnTrueWhenSpecifiedTimeIsDayAfterThanCurrentThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	now := time.Now()
	nextDay := now.Add(24 * time.Hour)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, nextDay)

	// then:
	assertions.NotNil(workoutGroup)
	assertions.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsEqualToThreshold_Unit(t *testing.T) {
	t.Log("This test should be implemented. Currently not found way to mock time in idomatic approach!")
	t.Skip()

	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"

	now := time.Now()
	threshold := 3 * time.Hour
	timeUnderLimit := now.Add(threshold)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, name, desc, timeUnderLimit)

	// then:
	assertions.NotNil(workoutGroup)
	assertions.Nil(err)
}
