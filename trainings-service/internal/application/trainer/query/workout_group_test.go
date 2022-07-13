package query_test

import (
	"context"
	"errors"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestWorkoutGroupHandler_ShouldReturnTrainerWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	service := mocks.NewTrainerService(t)
	SUT, _ := query.NewWorkoutGroupHandler(service)

	const (
		trainerUUID  = "467a73d6-3bf1-46b3-a977-02d0ab408787"
		customerUUID = "2e049075-4621-4a42-85e8-33b8d2c8ccef"
		customerName = "Jerry Doe"
	)

	group := newTestTrainerWorkoutGroup(trainerUUID)
	group.AssignCustomer(newTestCustomerDetails(customerUUID, customerName))
	groupUUID := group.UUID()
	expectedGroup := query.WorkoutGroupDetails{
		TrainerUUID: group.TrainerUUID(),
		TrainerName: group.TrainerName(),
		GroupUUID:   group.UUID(),
		GroupDesc:   group.Description(),
		GroupName:   group.Name(),
		Customers: []query.CustomerData{
			{
				UUID: customerUUID,
				Name: customerName,
			},
		},
		Date: group.Date().String(),
	}

	service.EXPECT().GetTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(group, nil)

	// when:
	actualGroup, err := SUT.Do(ctx, query.WorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedGroup, actualGroup)
	mock.AssertExpectationsForObjects(t, service)
}

func TestWorkoutGroupHandler_ShouldNotReturnTrainerWorkoutGroupWithSuccessWhenTrainerServiceFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	ctx := context.Background()
	service := mocks.NewTrainerService(t)
	SUT, _ := query.NewWorkoutGroupHandler(service)

	const (
		trainerUUID = "467a73d6-3bf1-46b3-a977-02d0ab408787"
		groupUUID   = "2e049075-4621-4a42-85e8-33b8d2c8ccef"
	)

	errServiceFailure := errors.New("trainer service failure")
	service.EXPECT().GetTrainerWorkoutGroup(ctx, trainerUUID, groupUUID).Return(trainer.WorkoutGroup{}, errServiceFailure)

	// when:
	actualGroup, err := SUT.Do(ctx, query.WorkoutGroupArgs{
		TrainerUUID: trainerUUID,
		GroupUUID:   groupUUID,
	})

	// then:
	assertions.ErrorIs(err, errServiceFailure)
	assertions.Empty(actualGroup)
	mock.AssertExpectationsForObjects(t, service)
}

func newTestCustomerDetails(customerUUID, name string) customer.Details {
	details, err := customer.NewCustomerDetails(customerUUID, name)
	if err != nil {
		panic(err)
	}
	return details
}
