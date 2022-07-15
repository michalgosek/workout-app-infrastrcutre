package query_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestWorkoutGroupsHandler_ShouldReturnTrainerWorkoutGroupDetailsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given
	const (
		trainerUUID  = "e41bc13b-9d4b-42ec-ab36-3bf6688b03fb"
		trainerName  = "John Doe"
		groupUUID    = "93d0cda3-cf47-48ac-bf3f-5a30ae04a327"
		groupDesc    = "dummy group desc"
		groupName    = "dummy group name"
		groupDate    = "12/03/2024 12:22"
		customerName = "Jerry Smith"
		customerUUID = "84f6f70e-8eb9-426a-894e-c074dfceae99"
	)
	ctx := context.Background()
	readModel := mocks.NewTrainerWorkoutGroupsReadModel(t)
	SUT := query.NewTrainerWorkoutGroupsHandler(readModel)

	expectedGroups := []query.TrainerWorkoutGroup{
		{
			TrainerUUID: trainerUUID,
			TrainerName: trainerName,
			GroupUUID:   groupUUID,
			GroupDesc:   groupDesc,
			GroupName:   groupName,
			Date:        groupDate,
			Participants: []query.TrainerWorkoutGroupParticipant{
				{
					UUID: customerUUID,
					Name: customerName,
				},
			},
		},
	}

	readModel.EXPECT().TrainerWorkoutGroups(ctx, trainerUUID).Return(expectedGroups, nil)

	// when:
	actualGroups, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroups, actualGroups)
	mock.AssertExpectationsForObjects(t, readModel)
}

func TestWorkoutGroupsHandler_ShouldReturnNotReturnTrainerWorkoutGroupDetailssWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given
	const trainerUUID = "e41bc13b-9d4b-42ec-ab36-3bf6688b03fb"
	ctx := context.Background()
	readModel := mocks.NewTrainerWorkoutGroupsReadModel(t)
	SUT := query.NewTrainerWorkoutGroupsHandler(readModel)

	trainerServiceFailureErr := errors.New("trainer service failure")
	readModel.EXPECT().TrainerWorkoutGroups(ctx, trainerUUID).Return(nil, trainerServiceFailureErr)

	// when:
	actualGroups, err := SUT.Do(ctx, trainerUUID)

	// then:
	assertions.ErrorIs(err, trainerServiceFailureErr)
	assertions.Empty(actualGroups)
	mock.AssertExpectationsForObjects(t, readModel)
}
