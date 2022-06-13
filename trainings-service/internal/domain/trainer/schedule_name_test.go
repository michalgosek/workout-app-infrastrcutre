package trainer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
)

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
