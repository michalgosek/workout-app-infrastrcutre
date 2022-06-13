package trainer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/adapters/testutil"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(customersAssigned, SUT.AssignedCustomers())
}

func AssignCustomerToTrainerSchedule(schedule *trainer.TrainerSchedule, n int) {
	for i := 0; i < n; i++ {
		schedule.AssignCustomer(uuid.NewString())
	}
}

func TestShouldNotReturnErrorWhenScheduleNameIsUnderLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	const newName = "dummy1"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assert.Nil(err)
	assert.Equal(newName, SUT.Name())

}

func TestShouldReturnErrorWhenScheduleNameIsOverLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	newName := strings.Repeat("s", 16)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assert.Equal(trainer.ErrScheduleNameExceeded, err)
}

func TestShouldNotReturnErrorWhenScheduleNameIsEqualLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	newName := strings.Repeat("s", 15)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assert.Nil(err)
}

func TestShouldRegisterCustomerToScheduleWithSucces(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 9
	SUT := testutil.GenerateTrainerSchedule(trainerUUID)

	// when:
	err := SUT.AssignCustomer(customerUUID)

	// then:
	assert.Nil(err)
	assert.Equal(customersLeft, SUT.Limit())
}

func TestShouldUnregisterCustomerFromScheduleWithSucces(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 10
	SUT := testutil.GenerateTrainerSchedule(trainerUUID)
	SUT.AssignCustomer(customerUUID)

	// when:
	SUT.UnregisterCustomer(customerUUID)

	// then:
	assert.Empty(SUT.AssignedCustomers())
	assert.Equal(customersLeft, SUT.Limit())
}

func TestShouldNotReturnErrorWhenTextLengthIsUnderLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	const newDesc = "dummy1"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(newDesc)

	// then:
	assert.Nil(err)
	assert.Equal(newDesc, SUT.Desc())
}

func TestShouldReturnErrorWhenTextLengthIsOverLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	invalidDesc := strings.Repeat("a", 101)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(invalidDesc)

	// then:
	assert.ErrorIs(err, trainer.ErrScheduleDescriptionExceeded)
	assert.Equal(desc, SUT.Desc())
}

func TestShouldReturnErrorWhenTextLengthEqualsLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	date := time.Now().Add(24 * time.Hour)
	expectedDesc := strings.Repeat("a", 100)
	SUT, _ := trainer.NewSchedule(trainerUUID, name, desc, date)

	// when:
	err := SUT.UpdateDesc(expectedDesc)

	// then:
	assert.Nil(err)
	assert.Equal(expectedDesc, SUT.Desc())
}

func TestShouldReturnErrorWhenSpecifiedTimeIsOneMinEarlierFromNow_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"

	threshold := 3 * time.Hour
	hourEarlier := time.Now().Add(threshold - 1*time.Hour)

	// when:
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, hourEarlier)

	// then:
	assert.Equal(trainer.ErrScheduleDateViolation, err)
	assert.Nil(schedule)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneMinLaterFromThreshold_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Hour)

	// when:
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, minLater)

	// then:
	assert.NotNil(schedule)
	assert.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneSecondLaterFromThreshold_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Second)

	// when:
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, minLater)

	// then:
	assert.NotNil(schedule)
	assert.Nil(err)
}

func TestShouldReturnTrueWhenSpecifiedTimeIsDayAfterThanCurrentThreshold_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"
	now := time.Now()
	nextDay := now.Add(24 * time.Hour)

	// when:
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, nextDay)

	// then:
	assert.NotNil(schedule)
	assert.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsEqualToThreshold_Unit(t *testing.T) {
	t.Log("This test should be implemented. Currently not found way to mock time in idomatic approach!")
	t.Skip()

	assert := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const name = "dummy"
	const desc = "dummy"

	now := time.Now()
	threshold := 3 * time.Hour
	timeUnderLimit := now.Add(threshold)

	// when:
	schedule, err := trainer.NewSchedule(trainerUUID, name, desc, timeUnderLimit)

	// then:
	assert.NotNil(schedule)
	assert.Nil(err)
}
