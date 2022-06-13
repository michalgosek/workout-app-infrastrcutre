package trainer_test

import (
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

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
