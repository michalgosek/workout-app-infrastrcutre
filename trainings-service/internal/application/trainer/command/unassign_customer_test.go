package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/command/mocks"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	"github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShouldUnassignCustomerToSpecifiedWorkoutGroupWithSuccess_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	_ = trainerWorkout.AssignCustomer(customerUUID)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	customerRepository := new(mocks.CustomerRepository)
	trainerRepository := new(mocks.TrainerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(nil)

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Nil(err)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenGroupNotOwnedByTrainer_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const firstTrainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"
	const secondTrainerUUID = "2877afdd-a15a-451c-8857-25075d626d2a"

	secondTrainerWorkout := testutil.NewTrainerWorkoutGroup(secondTrainerUUID)
	secondTrainerWorkout.AssignCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, secondTrainerWorkout.UUID(), secondTrainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, secondTrainerWorkout.UUID()).Return(secondTrainerWorkout, nil)
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  firstTrainerUUID,
	})

	// then:
	assertions.Equal(err, command.ErrWorkoutGroupNotOwner)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenQueryTrainerWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, errors.New("error"))
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenQueryCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(customer.WorkoutDay{}, errors.New("error"))

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenCustomerWorkoutDayNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	_ = trainerWorkout.AssignCustomer(customerUUID)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(customer.WorkoutDay{}, nil)
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Equal(err, command.ErrResourceNotFound)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWheDeleteCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(errors.New("error"))

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWheUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	_ = trainerWorkout.AssignCustomer(customerUUID)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(errors.New("error"))

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.WorkoutUnregister{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.ErrorIs(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}
