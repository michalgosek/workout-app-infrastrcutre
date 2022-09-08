package authorization_test

import (
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/authorization"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasScopeShouldReturnTrueForViewTrainerGroupsClaims_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	SUT := authorization.CustomClaims{Permissions: []string{
		string(authorization.ViewParticipantGroups),
		string(authorization.ViewTrainerGroups),
	}}

	// when:
	result := SUT.HasScope(authorization.ViewTrainerGroups)

	// then:
	assertions.True(result)
}
