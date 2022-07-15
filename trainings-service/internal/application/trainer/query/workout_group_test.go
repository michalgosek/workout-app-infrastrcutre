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

func TestWorkoutGroupHandler_ShouldReturnTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	readModel := mocks.NewTrainerWorkoutGroupReadModel(t)
	SUT := query.NewTrainerWorkoutGroupHandler(readModel)

	const (
		trainerUUID  = "467a73d6-3bf1-46b3-a977-02d0ab408787"
		trainerName  = "John Doe"
		groupUUID    = "2fe5b064-7b0f-4471-a70f-27e794eb8826"
		groupDesc    = "dummy group desc"
		groupName    = "dummy group name"
		groupDate    = "12/03/2024 12:33"
		customerUUID = "2e049075-4621-4a42-85e8-33b8d2c8ccef"
		customerName = "Jerry Doe"
	)

	expectedGroup := query.TrainerWorkoutGroup{
		TrainerUUID: trainerUUID,
		TrainerName: trainerName,
		GroupUUID:   groupUUID,
		GroupDesc:   groupDesc,
		GroupName:   groupName,
		Participants: []query.TrainerWorkoutGroupParticipant{
			{
				UUID: customerUUID,
				Name: customerName,
			},
		},
		Date: groupDate,
	}

	readModel.EXPECT().TrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(expectedGroup, nil)

	// when:
	actualGroup, err := SUT.Do(ctx, trainerUUID, groupUUID)

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
	mock.AssertExpectationsForObjects(t, readModel)
}

func TestWorkoutGroupHandler_ShouldNotReturnTrainerWorkoutGroupWithSuccessWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	readModel := mocks.NewTrainerWorkoutGroupReadModel(t)
	SUT := query.NewTrainerWorkoutGroupHandler(readModel)

	const (
		trainerUUID = "467a73d6-3bf1-46b3-a977-02d0ab408787"
		groupUUID   = "2e049075-4621-4a42-85e8-33b8d2c8ccef"
	)

	errServiceFailure := errors.New("trainer service failure")
	readModel.EXPECT().TrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(query.TrainerWorkoutGroup{}, errServiceFailure)

	// when:
	actualGroup, err := SUT.Do(ctx, trainerUUID, groupUUID)

	// then:
	assertions.ErrorIs(err, errServiceFailure)
	assertions.Empty(actualGroup)
	mock.AssertExpectationsForObjects(t, readModel)
}
