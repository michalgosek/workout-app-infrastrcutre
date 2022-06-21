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
	const customerName = "Jerry Smith"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	ctx := context.Background()
	customerRepository := new(mocks.CustomerRepository)
	trainerRepository := new(mocks.TrainerRepository)

	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(details)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(nil)

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
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
	const customerName = "Jerry Smith"
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const firstTrainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"
	const secondTrainerUUID = "2877afdd-a15a-451c-8857-25075d626d2a"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	secondTrainerWorkout := testutil.NewTrainerWorkoutGroup(secondTrainerUUID)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	secondTrainerWorkout.AssignCustomer(details)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, secondTrainerWorkout.UUID(), secondTrainerWorkout.Date())

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, secondTrainerWorkout.UUID()).Return(secondTrainerWorkout, nil)
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
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
	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, errors.New("error"))
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenQueryCustomerWorkoutDayFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(customer.WorkoutDay{}, errors.New("error"))
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWhenCustomerWorkoutDayNotExist_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerName = "Jerry Smith"
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)
	trainerWorkout.AssignCustomer(details)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(customer.WorkoutDay{}, nil)
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
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

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(errors.New("error"))
	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	customerRepository.AssertExpectations(t)
	trainerRepository.AssertExpectations(t)
}

func TestShouldReturnErrorWheUpsertWorkoutGroupFailure_Unit(t *testing.T) {
	assertions := assert.New(t)

	// given:
	const customerName = "Jerry Smith"
	const customerUUID = "094bb50a-7da3-461f-86f6-46d16c055e1e"
	const trainerUUID = "090d4e58-3a5e-4eaf-8905-14b892d35678"

	ctx := context.Background()
	trainerRepository := new(mocks.TrainerRepository)
	customerRepository := new(mocks.CustomerRepository)
	trainerWorkout := testutil.NewTrainerWorkoutGroup(trainerUUID)
	details, _ := customer.NewCustomerDetails(customerUUID, customerName)

	_ = trainerWorkout.AssignCustomer(details)
	trainerWorkoutWithoutCustomer := trainerWorkout
	trainerWorkoutWithoutCustomer.UnregisterCustomer(customerUUID)
	customerWorkout, _ := customer.NewWorkoutDay(customerUUID, trainerWorkout.UUID(), trainerWorkout.Date())

	customerRepository.EXPECT().QueryCustomerWorkoutDay(ctx, customerUUID, trainerWorkout.UUID()).Return(*customerWorkout, nil)
	trainerRepository.EXPECT().QueryTrainerWorkoutGroup(ctx, trainerWorkout.UUID()).Return(trainerWorkout, nil)
	customerRepository.EXPECT().DeleteCustomerWorkoutDay(ctx, customerUUID, customerWorkout.UUID()).Return(nil)
	trainerRepository.EXPECT().UpsertTrainerWorkoutGroup(ctx, trainerWorkoutWithoutCustomer).Return(errors.New("error"))

	SUT := command.NewUnassignCustomerHandler(customerRepository, trainerRepository)

	// when:
	err := SUT.Do(ctx, command.UnassignCustomer{
		CustomerUUID: customerWorkout.CustomerUUID(),
		GroupUUID:    customerWorkout.GroupUUID(),
		TrainerUUID:  trainerWorkout.TrainerUUID(),
	})

	// then:
	assertions.Equal(err, command.ErrRepositoryFailure)
	trainerRepository.AssertExpectations(t)
	customerRepository.AssertExpectations(t)
}
