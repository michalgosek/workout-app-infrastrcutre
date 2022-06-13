package trainer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

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
