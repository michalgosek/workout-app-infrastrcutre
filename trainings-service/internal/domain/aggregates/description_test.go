package aggregates_test

import (
	"strings"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/aggregates"
	"github.com/stretchr/testify/assert"
)

func TesShouldNotReturnErrorWhenTextLengthIsUnderLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	s := strings.Repeat("a", 1)
	const limit = 3
	SUT := aggregates.NewWorkoutDescription(limit)

	// when:
	err := SUT.Check(s)

	// then:
	assert.Nil(err)
}

func TestShouldReturnErrorWhenTextLengthIsOverLimit_Unit(t *testing.T) {
	assert := assert.New(t)

	// given:
	const limit = 3
	s := strings.Repeat("a", 100)
	SUT := aggregates.NewWorkoutDescription(limit)

	// when:
	err := SUT.Check(s)

	// then:
	assert.ErrorIs(err, aggregates.ErrWorkoutSessionDescriptionExceeded)
}
