package aggregates_test

import (
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/aggregates"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenSpecifiedTimeIsOneMinEarlierFromNow_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const limit = 3
	threshold := limit * time.Hour
	hourEarlier := time.Now().Add(threshold - 1*time.Hour)
	SUT := aggregates.NewWorkoutDate(limit)

	// when:
	err := SUT.Check(hourEarlier)

	// then:
	assert.ErrorIs(err, aggregates.ErrDateAggregateViolation)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneMinLaterFromCurrentAggregateLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const limit = 3
	threshold := limit * time.Hour
	minLater := time.Now().Add(threshold + time.Hour)
	SUT := aggregates.NewWorkoutDate(limit)

	// when:
	err := SUT.Check(minLater)

	// then:
	assert.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsOneSecondLaterFromCurrentAggregateLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const limit = 3
	threshold := limit * time.Hour
	minLater := time.Now().Add(threshold + time.Second)
	SUT := aggregates.NewWorkoutDate(limit)

	// when:
	err := SUT.Check(minLater)

	// then:
	assert.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsDayAfterThanCurrenDay_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const limit = 3
	SUT := aggregates.NewWorkoutDate(limit)
	now := time.Now()
	nextDay := now.Add(24 * time.Hour)

	// when:
	err := SUT.Check(nextDay)

	// then:
	assert.Nil(err)
}

func TestShouldNotReturnErrorWhenSpecifiedTimeIsEqualToAggregateLimit_Unit(t *testing.T) {
	t.Log("This test should be implemented. Currently not found way to mock time in idomatic approach!")
	t.Skip()

	assert := assert.New(t)

	// given:
	const limit = 3
	SUT := aggregates.NewWorkoutDate(limit)
	now := time.Now()
	timeUnderLimit := now.Add(limit * time.Hour)

	// when:
	err := SUT.Check(timeUnderLimit)

	// then:
	assert.Nil(err)
}
