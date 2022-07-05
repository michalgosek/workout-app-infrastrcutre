// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	customer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/customer"
	mock "github.com/stretchr/testify/mock"
)

// CustomerRepository is an autogenerated mock type for the CustomerRepository type
type CustomerRepository struct {
	mock.Mock
}

type CustomerRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *CustomerRepository) EXPECT() *CustomerRepository_Expecter {
	return &CustomerRepository_Expecter{mock: &_m.Mock}
}

// DeleteCustomerWorkoutDay provides a mock function with given fields: ctx, customerUUID, groupUUID
func (_m *CustomerRepository) DeleteCustomerWorkoutDay(ctx context.Context, customerUUID string, groupUUID string) error {
	ret := _m.Called(ctx, customerUUID, groupUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, customerUUID, groupUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CustomerRepository_DeleteCustomerWorkoutDay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCustomerWorkoutDay'
type CustomerRepository_DeleteCustomerWorkoutDay_Call struct {
	*mock.Call
}

// DeleteCustomerWorkoutDay is a helper method to define mock.On call
//  - ctx context.Context
//  - customerUUID string
//  - groupUUID string
func (_e *CustomerRepository_Expecter) DeleteCustomerWorkoutDay(ctx interface{}, customerUUID interface{}, groupUUID interface{}) *CustomerRepository_DeleteCustomerWorkoutDay_Call {
	return &CustomerRepository_DeleteCustomerWorkoutDay_Call{Call: _e.mock.On("DeleteCustomerWorkoutDay", ctx, customerUUID, groupUUID)}
}

func (_c *CustomerRepository_DeleteCustomerWorkoutDay_Call) Run(run func(ctx context.Context, customerUUID string, groupUUID string)) *CustomerRepository_DeleteCustomerWorkoutDay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CustomerRepository_DeleteCustomerWorkoutDay_Call) Return(_a0 error) *CustomerRepository_DeleteCustomerWorkoutDay_Call {
	_c.Call.Return(_a0)
	return _c
}

// QueryCustomerWorkoutDay provides a mock function with given fields: ctx, customerUUID, GroupUUID
func (_m *CustomerRepository) QueryCustomerWorkoutDay(ctx context.Context, customerUUID string, GroupUUID string) (customer.WorkoutDay, error) {
	ret := _m.Called(ctx, customerUUID, GroupUUID)

	var r0 customer.WorkoutDay
	if rf, ok := ret.Get(0).(func(context.Context, string, string) customer.WorkoutDay); ok {
		r0 = rf(ctx, customerUUID, GroupUUID)
	} else {
		r0 = ret.Get(0).(customer.WorkoutDay)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, customerUUID, GroupUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CustomerRepository_QueryCustomerWorkoutDay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryCustomerWorkoutDay'
type CustomerRepository_QueryCustomerWorkoutDay_Call struct {
	*mock.Call
}

// QueryCustomerWorkoutDay is a helper method to define mock.On call
//  - ctx context.Context
//  - customerUUID string
//  - GroupUUID string
func (_e *CustomerRepository_Expecter) QueryCustomerWorkoutDay(ctx interface{}, customerUUID interface{}, GroupUUID interface{}) *CustomerRepository_QueryCustomerWorkoutDay_Call {
	return &CustomerRepository_QueryCustomerWorkoutDay_Call{Call: _e.mock.On("QueryCustomerWorkoutDay", ctx, customerUUID, GroupUUID)}
}

func (_c *CustomerRepository_QueryCustomerWorkoutDay_Call) Run(run func(ctx context.Context, customerUUID string, GroupUUID string)) *CustomerRepository_QueryCustomerWorkoutDay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *CustomerRepository_QueryCustomerWorkoutDay_Call) Return(_a0 customer.WorkoutDay, _a1 error) *CustomerRepository_QueryCustomerWorkoutDay_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// UpsertCustomerWorkoutDay provides a mock function with given fields: ctx, workout
func (_m *CustomerRepository) UpsertCustomerWorkoutDay(ctx context.Context, workout customer.WorkoutDay) error {
	ret := _m.Called(ctx, workout)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, customer.WorkoutDay) error); ok {
		r0 = rf(ctx, workout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CustomerRepository_UpsertCustomerWorkoutDay_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertCustomerWorkoutDay'
type CustomerRepository_UpsertCustomerWorkoutDay_Call struct {
	*mock.Call
}

// UpsertCustomerWorkoutDay is a helper method to define mock.On call
//  - ctx context.Context
//  - workout customer.WorkoutDay
func (_e *CustomerRepository_Expecter) UpsertCustomerWorkoutDay(ctx interface{}, workout interface{}) *CustomerRepository_UpsertCustomerWorkoutDay_Call {
	return &CustomerRepository_UpsertCustomerWorkoutDay_Call{Call: _e.mock.On("UpsertCustomerWorkoutDay", ctx, workout)}
}

func (_c *CustomerRepository_UpsertCustomerWorkoutDay_Call) Run(run func(ctx context.Context, workout customer.WorkoutDay)) *CustomerRepository_UpsertCustomerWorkoutDay_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(customer.WorkoutDay))
	})
	return _c
}

func (_c *CustomerRepository_UpsertCustomerWorkoutDay_Call) Return(_a0 error) *CustomerRepository_UpsertCustomerWorkoutDay_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewCustomerRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCustomerRepository creates a new instance of CustomerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCustomerRepository(t mockConstructorTestingTNewCustomerRepository) *CustomerRepository {
	mock := &CustomerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
