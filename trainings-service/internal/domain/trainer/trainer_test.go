package trainer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenCustomerLimitExceeded_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const trainerUUID = "346dcf15-549f-4853-aa92-6ecbc6486ce8"
	const customerUUID = "5b6bd420-2b8a-444f-869a-ea12957ef8c1"
	const customersLeft = 0
	const customersAssigned = 10

	SUT := testutil.NewTrainerWorkoutGroup(trainerUUID)
	AssignCustomerToWorkoutGroup(&SUT, 10)

	// when:
	err := SUT.AssignCustomer(customerUUID)

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
	const newName = "dummy1"
	const trainerName = "John Doe"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateGroupName(newName)

	// then:
	assertions.Nil(err)
	assertions.Equal(newName, SUT.GroupName())

}

func TestShouldReturnErrorWhenWorkoutGroupNameIsOverLimit_Unit(t *testing.T) {
	assertions := assert.New(t)

	const trainerUUID = "1b0af14e-5aa9-4b80-968f-03d93f46805e"
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	newName := strings.Repeat("s", 16)
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateGroupName(newName)

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
	err := SUT.UpdateGroupName(newName)

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
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.AssignCustomer(customerUUID)

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
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

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
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"
	const newDesc = "dummy1"
	date := time.Now().Add(24 * time.Hour)
	SUT, _ := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, date)

	// when:
	err := SUT.UpdateGroupDescription(newDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(newDesc, SUT.GroupDescription())
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
	err := SUT.UpdateGroupDescription(invalidDesc)

	// then:
	assertions.ErrorIs(err, trainer.ErrScheduleDescriptionExceeded)
	assertions.Equal(groupDesc, SUT.GroupDescription())
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
	err := SUT.UpdateGroupDescription(expectedDesc)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedDesc, SUT.GroupDescription())
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
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, hourEarlier)

	// then:
	assertions.Equal(trainer.ErrScheduleDateViolation, err)
	assertions.Nil(workoutGroup)
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
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.NotNil(workoutGroup)
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
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.NotNil(workoutGroup)
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
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, nextDay)

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
	const groupName = "dummy"
	const groupDesc = "dummy"
	const trainerName = "John Doe"

	now := time.Now()
	threshold := 3 * time.Hour
	timeUnderLimit := now.Add(threshold)

	// when:
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, timeUnderLimit)

	// then:
	assertions.NotNil(workoutGroup)
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
	date := time.Now().Add(24 * time.Hour)

	expectedWorkoutGroup, _ := trainer.NewWorkoutGroup(groupUUID, trainerName, groupName, groupDesc, date)
	expectedWorkoutGroup.AssignCustomer(customerUUID)

	// when:
	workoutGroup, err := trainer.UnmarshalFromDatabase(
		expectedWorkoutGroup.UUID(),
		expectedWorkoutGroup.TrainerUUID(),
		expectedWorkoutGroup.TrainerName(),
		expectedWorkoutGroup.GroupName(),
		expectedWorkoutGroup.GroupDescription(),
		expectedWorkoutGroup.CustomerUUIDs(),
		expectedWorkoutGroup.Date(),
		expectedWorkoutGroup.Limit(),
	)

	// then:
	assertions.Nil(err)
	assertions.Equal(*expectedWorkoutGroup, workoutGroup)
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
	workoutGroup, err := trainer.NewWorkoutGroup(trainerUUID, trainerName, groupName, groupDesc, minLater)

	// then:
	assertions.Equal(trainer.ErrEmptyTrainerName, err)
	assertions.Nil(workoutGroup)
}

func AssignCustomerToWorkoutGroup(workoutGroup *trainer.WorkoutGroup, n int) {
	for i := 0; i < n; i++ {
		err := workoutGroup.AssignCustomer(uuid.NewString())
		if err != nil {
			panic(err)
		}
	}
}
