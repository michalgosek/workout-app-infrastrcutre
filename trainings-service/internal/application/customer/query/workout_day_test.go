package query_test

import (
	"context"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/query"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/customer/query/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldGetRequestedWorkoutDayWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	const customerName = "John Doe"
	const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"

	ctx := context.Background()
	repository := new(mocks.WorkoutDayHandlerRepository)
	SUT := query.NewWorkoutDayHandler(repository)

	customerWorkoutDay := testutil.NewWorkoutDay(customerUUID)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerDetails, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(customerDetails)

	repository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.TrainerUUID(), trainerWorkout.UUID()).Return(trainerWorkout, nil)
	repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(customerWorkoutDay, nil)

	expectedDay := query.CustomerWorkoutDay{
		Date:         trainerWorkout.Date(),
		Trainer:      trainerWorkout.TrainerName(),
		WorkoutName:  trainerWorkout.Name(),
		WorkoutDesc:  trainerWorkout.Description(),
		Participants: trainerWorkout.AssignedCustomers(),
	}

	// when:
	day, err := SUT.Do(ctx, query.WorkoutDay{
		CustomerUUID: customerUUID,
		GroupUUID:    trainerWorkout.UUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Nil(err)
	assertions.Equal(expectedDay, day)
	repository.AssertExpectations(t)
}
