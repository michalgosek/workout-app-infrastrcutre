package trainer_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenCustomerLimitExceeded_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const customersLeft = 0
	const customerName = "John Doe"
	const customersAssigned = 10
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	date := time.Now().Add(24 * time.Hour)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)
	_ = AssignCustomerToWorkoutGroup(&SUT, 10)

	// when:
	err := SUT.AssignCustomer(details)

	// then:
	assertions.ErrorIs(trainer.ErrCustomersScheduleLimitExceeded, err)
	assertions.Equal(customersLeft, SUT.Limit())
	assertions.Equal(customersAssigned, SUT.AssignedCustomers())
}

func TestShouldNotReturnErrorWhenWorkoutGroupNameIsUnderLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const updateName = "dummy1"
	const trainerName = "John Doe"

	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateName(updateName)

	// then:
	assertions.Nil(err)
	assertions.Equal(updateName, SUT.Name())
}

func TestShouldReturnErrorWhenWorkoutGroupNameIsOverLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	updatedDesc := strings.Repeat("s", 16)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateName(updatedDesc)

	// then:
	assertions.Equal(trainer.ErrScheduleNameExceeded, err)
}

func TestShouldNotReturnErrorWhenWorkoutGroupNameIsEqualLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	newName := strings.Repeat("s", 15)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, groupName, groupDesc, trainerName, date)

	// when:
	err := SUT.UpdateName(newName)

	// then:
	assertions.Nil(err)
}

func TestShouldRegisterCustomerToWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 9
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	const customerName = "Jerry Smith"

	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)

	// when:
	err := SUT.AssignCustomer(details)

	// then:
	assertions.Nil(err)
	assertions.Equal(customersLeft, SUT.Limit())
}

func TestShouldUnregisterCustomerFromWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const customerUUID = "1b83c88b-4aac-4719-ac23-03a43627cb3e"
	const customersLeft = 10
	const customerName = "Jerry Smith"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	_ = SUT.AssignCustomer(details)

	// when:
	SUT.UnregisterCustomer(details.UUID())

	// then:
	assertions.Empty(SUT.AssignedCustomers())
	assertions.Equal(customersLeft, SUT.Limit())
}

func TestShouldNotReturnErrorWhenTextLengthIsUnderLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	const newDesc = "dummy1"

	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateDescription(newDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(newDesc, SUT.Description())
}

func TestShouldReturnErrorWhenTextLengthIsOverLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	date := time.Now().Add(24 * time.Hour)
	invalidDesc := strings.Repeat("a", 101)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateDescription(invalidDesc)

	// then:
	assertions.ErrorIs(err, trainer.ErrScheduleDescriptionExceeded)
	assertions.Equal(groupDesc, SUT.Description())
}

func TestShouldReturnErrorWhenTextLengthEqualsLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	date := time.Now().Add(24 * time.Hour)
	expectedDesc := strings.Repeat("a", 100)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateDescription(expectedDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedDesc, SUT.Description())
}

func TestShouldReturnErrorWhenSpecifiedTimeIsOneMinEarlierFromNow_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	threshold := 3 * time.Hour
	hourEarlier := time.Now().Add(threshold - 1*time.Hour)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, hourEarlier)

	// then:
	assertions.Equal(trainer.ErrScheduleDateViolation, err)
	assertions.Empty(SUT)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneMinLaterFromThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Hour)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.NotNil(SUT)
	assertions.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneSecondLaterFromThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	threshold := 3 * time.Hour
	minLater := time.Now().Add(threshold + time.Second)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.NotNil(SUT)
	assertions.Nil(err)
}

func TestShouldReturnTrueWhenSpecifiedTimeIsDayAfterThanCurrentThreshold_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	now := time.Now()
	nextDay := now.Add(24 * time.Hour)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, nextDay)

	// then:
	assertions.NotNil(SUT)
	assertions.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsEqualToThreshold_Unit(t *testing.T) {
	t.Log("This test should be implemented. Currently not found way to mock time in idomatic approach!")
	t.Skip()

	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	now := time.Now()
	threshold := 3 * time.Hour
	timeUnderLimit := now.Add(threshold)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, timeUnderLimit)

	// then:
	assertions.NotNil(SUT)
	assertions.Nil(err)
}

func TestUnmarshalFromDatabaseShouldParseDataWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const groupUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	const customerUUID = "fb561c94-c60a-4864-84cb-9901cabf9ed5"
	const customerName = "Jerry Smith"
	date := time.Now().Add(24 * time.Hour)

	expectedWorkoutGroup, _ := trainer.NewWorkoutGroup(groupUUID, trainerName, groupName, groupDesc, date)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	expectedWorkoutGroup.AssignCustomer(details)

	// when:
	SUT, err := trainer.UnmarshalWorkoutGroupFromDatabase(
		trainer.WorkoutGroupDetails{
			UUID:        expectedWorkoutGroup.UUID(),
			TrainerUUID: expectedWorkoutGroup.TrainerUUID(),
			TrainerName: expectedWorkoutGroup.TrainerName(),
			Name:        expectedWorkoutGroup.Name(),
			Description: expectedWorkoutGroup.Description(),
			Date:        expectedWorkoutGroup.Date(),
			Limit:       expectedWorkoutGroup.Limit(),
		},
		expectedWorkoutGroup.CustomerDetails())

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkoutGroup, SUT)
}

func TestShouldNotReturnErrorWhenSpecifiedTrainerNameIsEmpty_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = ""

	minLater := time.Now().Add(24 * time.Hour)

	// when:
	SUT, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.Equal(trainer.ErrEmptyTrainerName, err)
	assertions.Empty(SUT)
}

func AssignCustomerToWorkoutGroup(workoutGroup *trainer.WorkoutGroup, n int) error {
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("Person%d", i+1)
		details, err := customer.NewCustomerDetails(uuid.NewString(), name)
		if err != nil {
			return err
		}
		err = workoutGroup.AssignCustomer(details)
		if err != nil {
			return err
		}
	}
	return nil
}
