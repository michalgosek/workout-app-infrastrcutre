package query_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestWorkoutGroupsHandler_ShouldReturnTrainerWorkoutGroupDetailsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given
	const trainerUUID = "e41bc13b-9d4b-42ec-ab36-3bf6688b03fb"
	ctx := context.Background()
	service := mocks.NewTrainerService(t)
	SUT, _ := query.NewWorkoutGroupsHandler(service)

	groups := []trainer.WorkoutGroup{
		newTestTrainerWorkoutGroup(trainerUUID),
		newTestTrainerWorkoutGroup(trainerUUID),
	}
	expectedGroupDetails := query.WorkoutGroupsDetails{
		WorkoutGroups: query.ConvertToWorkoutGroupsDetails(groups...),
	}

	service.EXPECT().GetTrainerWorkoutGroups(ctx, trainerUUID).Return(groups, nil)

	// when:
	actualGroups, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroupDetails, actualGroups)
	mock.AssertExpectationsForObjects(t, service)
}

func TestWorkoutGroupsHandler_ShouldReturnNotReturnTrainerWorkoutGroupDetailssWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given
	const trainerUUID = "e41bc13b-9d4b-42ec-ab36-3bf6688b03fb"
	ctx := context.Background()
	service := mocks.NewTrainerService(t)
	SUT, _ := query.NewWorkoutGroupsHandler(service)

	trainerServiceFailureErr := errors.New("trainer service failure")
	service.EXPECT().GetTrainerWorkoutGroups(ctx, trainerUUID).Return(nil, trainerServiceFailureErr)

	// when:
	actualGroups, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.ErrorIs(err, trainerServiceFailureErr)
	assertions.Empty(actualGroups)
	mock.AssertExpectationsForObjects(t, service)
}

func newTestTrainerWorkoutGroup(trainerUUID string) trainer.WorkoutGroup {
	schedule := time.Now().AddDate(0, 0, 1)
	group, err := trainer.NewWorkoutGroup(trainerUUID, "dummy_trainer", "dummy_group", "dummy_desc", schedule)
	if err != nil {
		panic(err)
	}
	return group
}
