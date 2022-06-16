// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	trainer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	mock "github.com/stretchr/testify/mock"
)

// WorkoutDeleter is an autogenerated mock type for the WorkoutDeleter type
type WorkoutDeleter struct {
	mock.Mock
}

type WorkoutDeleter_Expecter struct {
	mock *mock.Mock
}

func (_m *WorkoutDeleter) EXPECT() *WorkoutDeleter_Expecter {
	return &WorkoutDeleter_Expecter{mock: &_m.Mock}
}

// DeleteWorkoutGroup provides a mock function with given fields: ctx, groupUUID
func (_m *WorkoutDeleter) DeleteWorkoutGroup(ctx context.Context, groupUUID string) error {
	ret := _m.Called(ctx, groupUUID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, groupUUID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WorkoutDeleter_DeleteWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteWorkoutGroup'
type WorkoutDeleter_DeleteWorkoutGroup_Call struct {
	*mock.Call
}

// DeleteWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - groupUUID string
func (_e *WorkoutDeleter_Expecter) DeleteWorkoutGroup(ctx interface{}, groupUUID interface{}) *WorkoutDeleter_DeleteWorkoutGroup_Call {
	return &WorkoutDeleter_DeleteWorkoutGroup_Call{Call: _e.mock.On("DeleteWorkoutGroup", ctx, groupUUID)}
}

func (_c *WorkoutDeleter_DeleteWorkoutGroup_Call) Run(run func(ctx context.Context, groupUUID string)) *WorkoutDeleter_DeleteWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *WorkoutDeleter_DeleteWorkoutGroup_Call) Return(_a0 error) *WorkoutDeleter_DeleteWorkoutGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

// QueryWorkoutGroup provides a mock function with given fields: ctx, groupUUID
func (_m *WorkoutDeleter) QueryWorkoutGroup(ctx context.Context, groupUUID string) (trainer.WorkoutGroup, error) {
	ret := _m.Called(ctx, groupUUID)

	var r0 trainer.WorkoutGroup
	if rf, ok := ret.Get(0).(func(context.Context, string) trainer.WorkoutGroup); ok {
		r0 = rf(ctx, groupUUID)
	} else {
		r0 = ret.Get(0).(trainer.WorkoutGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, groupUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WorkoutDeleter_QueryWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryWorkoutGroup'
type WorkoutDeleter_QueryWorkoutGroup_Call struct {
	*mock.Call
}

// QueryWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - groupUUID string
func (_e *WorkoutDeleter_Expecter) QueryWorkoutGroup(ctx interface{}, groupUUID interface{}) *WorkoutDeleter_QueryWorkoutGroup_Call {
	return &WorkoutDeleter_QueryWorkoutGroup_Call{Call: _e.mock.On("QueryWorkoutGroup", ctx, groupUUID)}
}

func (_c *WorkoutDeleter_QueryWorkoutGroup_Call) Run(run func(ctx context.Context, groupUUID string)) *WorkoutDeleter_QueryWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *WorkoutDeleter_QueryWorkoutGroup_Call) Return(_a0 trainer.WorkoutGroup, _a1 error) *WorkoutDeleter_QueryWorkoutGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}
