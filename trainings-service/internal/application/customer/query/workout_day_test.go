package query_test

import (
	"testing"
)

func TestShouldGetRequestedWorkoutDayWithSuccess_Unit(t *testing.T) {
	//assertions := assert.New(t)
	//
	//// given:
	//ctx := context.Background()
	//const customerUUID = "f2691a1e-575e-4fa8-8a37-e01d29a204e1"
	//const trainerUUID = "c7ea5361-faec-4d69-9eff-86c3e10384a9"
	//repository := new(mocks.WorkoutDayHandlerRepository)
	//SUT := query.NewWorkoutDayHandler(repository)
	//
	//customerWorkoutDay := testutil.NewWorkoutDay(customerUUID)
	//trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	//trainerWorkout.AssignCustomer(customerUUID)
	//
	//repository.EXPECT().QueryWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	//repository.EXPECT().QueryCustomerWorkoutDay(ctx, customerWorkoutDay.UUID(), customerWorkoutDay.CustomerUUID()).Return(customerWorkoutDay, nil)
	//
	//expectedDay := query.CustomerWorkoutDay{
	//	Date:         trainerWorkout.Date(),
	//	Trainer:      trainerWorkout.TrainerName(),
	//	WorkoutName:  trainerWorkout.GroupName(),
	//	WorkoutDesc:  trainerWorkout.GroupDescription(),
	//	Participants: trainerWorkout.Limit(),
	//}
	//
	//// when:
	//day, err := SUT.Do(ctx, query.WorkoutDayDetails{
	//	CustomerUUID: customerUUID,
	//	GroupUUID:    trainerWorkout.UUID(),
	//})
	//// then:
	//assertions.Nil(err)
	//assertions.Equal(expectedDay, day)
	//repository.AssertExpectations(t)
}
