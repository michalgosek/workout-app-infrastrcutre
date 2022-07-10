package query_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/query/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainer"
	domain "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestWorkoutDayHandler_ShouldReturnCustomerWorkoutGroupDetailsWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "1ef5dddb-f809-427a-b3d7-fb758c883ed1"
		customerUUID = "577f3ca5-4302-4509-8820-9899a82e054c"
	)

	trainerService := new(mocks.TrainerService)
	SUT, _ := query.NewWorkoutDayHandler(trainerService)

	ctx := context.Background()
	trainerWorkoutGroupWithCustomer := newTestTrainerWorkoutGroup(trainerUUID)
	expectedWorkoutGroupDetails := query.CustomerWorkoutGroupDetails{
		TrainerUUID:  trainerUUID,
		TrainerName:  trainerWorkoutGroupWithCustomer.TrainerName(),
		WorkoutName:  trainerWorkoutGroupWithCustomer.Name(),
		WorkoutDesc:  trainerWorkoutGroupWithCustomer.Description(),
		Date:         trainerWorkoutGroupWithCustomer.Date(),
		Participants: trainerWorkoutGroupWithCustomer.AssignedCustomers(),
	}
	groupUUID := trainerWorkoutGroupWithCustomer.UUID()

	trainerService.EXPECT().GetCustomerWorkoutGroup(ctx, trainer.WorkoutGroupWithCustomerArgs{
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
		CustomerUUID: customerUUID,
	}).Return(trainerWorkoutGroupWithCustomer, nil)

	// when:
	actualCustomerWorkoutGroupDetails, err := SUT.Do(ctx, query.WorkoutDayHandlerArgs{
		CustomerUUID: customerUUID,
		TrainerUUID:  trainerUUID,
		GroupUUID:    groupUUID,
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedWorkoutGroupDetails, actualCustomerWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, trainerService)
}

func TestWorkoutDayHandler_ShouldReturnNotCustomerWorkoutGroupDetailsWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const (
		trainerUUID  = "1ef5dddb-f809-427a-b3d7-fb758c883ed1"
		customerUUID = "577f3ca5-4302-4509-8820-9899a82e054c"
		groupUUD     = "0b68287d-76cd-407f-a644-c076b59b5645"
	)

	trainerService := new(mocks.TrainerService)
	SUT, _ := query.NewWorkoutDayHandler(trainerService)

	ctx := context.Background()
	repositoryFailureErr := errors.New("repository failure")
	trainerService.EXPECT().GetCustomerWorkoutGroup(ctx, trainer.WorkoutGroupWithCustomerArgs{
		TrainerUUID:  trainerUUID,
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUD,
	}).Return(domain.WorkoutGroup{}, repositoryFailureErr)

	// when:
	actualCustomerWorkoutGroupDetails, err := SUT.Do(ctx, query.WorkoutDayHandlerArgs{
		CustomerUUID: customerUUID,
		GroupUUID:    groupUUD,
		TrainerUUID:  trainerUUID,
	})

	// then:
	assertions.ErrorIs(err, repositoryFailureErr)
	assertions.Empty(actualCustomerWorkoutGroupDetails)
	mock.AssertExpectationsForObjects(t, trainerService)
}

func newTestTrainerWorkoutGroup(trainerUUID string) domain.WorkoutGroup {
	schedule := time.Now().AddDate(0, 0, 1)
	group, err := domain.NewWorkoutGroup(trainerUUID, "dummy_trainer", "dummy_group", "dummy_desc", schedule)
	if err != nil {
		panic(err)
	}
	return group
}
