package http_test

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/ports/http"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldReturnErrorBadRequestForMissingCustomerUUIDWhenCancelWorkout_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:

	SUT := http.NewCustomerHTTP()

	// when:

	// then:

}
