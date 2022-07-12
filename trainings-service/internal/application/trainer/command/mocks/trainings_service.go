// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	trainings "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/services/trainings"
	mock "github.com/stretchr/testify/mock"
)

// TrainingsService is an autogenerated mock type for the TrainingsService type
type TrainingsService struct {
	mock.Mock
}

type TrainingsService_Expecter struct {
	mock *mock.Mock
}

func (_m *TrainingsService) EXPECT() *TrainingsService_Expecter {
	return &TrainingsService_Expecter{mock: &_m.Mock}
}

// CancelTrainerWorkoutGroup provides a mock function with given fields: ctx, args
func (_m *TrainingsService) CancelTrainerWorkoutGroup(ctx context.Context, args trainings.CancelTrainerWorkoutGroupArgs) error {
	ret := _m.Called(ctx, args)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, trainings.CancelTrainerWorkoutGroupArgs) error); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TrainingsService_CancelTrainerWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CancelTrainerWorkoutGroup'
type TrainingsService_CancelTrainerWorkoutGroup_Call struct {
	*mock.Call
}

// CancelTrainerWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - args trainings.CancelTrainerWorkoutGroupArgs
func (_e *TrainingsService_Expecter) CancelTrainerWorkoutGroup(ctx interface{}, args interface{}) *TrainingsService_CancelTrainerWorkoutGroup_Call {
	return &TrainingsService_CancelTrainerWorkoutGroup_Call{Call: _e.mock.On("CancelTrainerWorkoutGroup", ctx, args)}
}

func (_c *TrainingsService_CancelTrainerWorkoutGroup_Call) Run(run func(ctx context.Context, args trainings.CancelTrainerWorkoutGroupArgs)) *TrainingsService_CancelTrainerWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(trainings.CancelTrainerWorkoutGroupArgs))
	})
	return _c
}

func (_c *TrainingsService_CancelTrainerWorkoutGroup_Call) Return(_a0 error) *TrainingsService_CancelTrainerWorkoutGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

// CancelTrainerWorkoutGroups provides a mock function with given fields: ctx, trainerUUID
func (_m *TrainingsService) CancelTrainerWorkoutGroups(ctx context.Context, trainerUUID string) error {
	ret := _m.Called(ctx, trainerUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, trainerUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TrainingsService_CancelTrainerWorkoutGroups_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CancelTrainerWorkoutGroups'
type TrainingsService_CancelTrainerWorkoutGroups_Call struct {
	*mock.Call
}

// CancelTrainerWorkoutGroups is a helper method to define mock.On call
//  - ctx context.Context
//  - trainerUUID string
func (_e *TrainingsService_Expecter) CancelTrainerWorkoutGroups(ctx interface{}, trainerUUID interface{}) *TrainingsService_CancelTrainerWorkoutGroups_Call {
	return &TrainingsService_CancelTrainerWorkoutGroups_Call{Call: _e.mock.On("CancelTrainerWorkoutGroups", ctx, trainerUUID)}
}

func (_c *TrainingsService_CancelTrainerWorkoutGroups_Call) Run(run func(ctx context.Context, trainerUUID string)) *TrainingsService_CancelTrainerWorkoutGroups_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *TrainingsService_CancelTrainerWorkoutGroups_Call) Return(_a0 error) *TrainingsService_CancelTrainerWorkoutGroups_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewTrainingsService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTrainingsService creates a new instance of TrainingsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTrainingsService(t mockConstructorTestingTNewTrainingsService) *TrainingsService {
	mock := &TrainingsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
