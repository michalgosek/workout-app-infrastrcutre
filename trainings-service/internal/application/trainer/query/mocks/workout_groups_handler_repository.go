// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	trainer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
)

// WorkoutGroupsHandlerRepository is an autogenerated mock type for the WorkoutGroupsHandlerRepository type
type WorkoutGroupsHandlerRepository struct {
	mock.Mock
}

type WorkoutGroupsHandlerRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *WorkoutGroupsHandlerRepository) EXPECT() *WorkoutGroupsHandlerRepository_Expecter {
	return &WorkoutGroupsHandlerRepository_Expecter{mock: &_m.Mock}
}

// QueryWorkoutGroups provides a mock function with given fields: ctx, trainerUUID
func (_m *WorkoutGroupsHandlerRepository) QueryWorkoutGroups(ctx context.Context, trainerUUID string) ([]trainer.WorkoutGroup, error) {
	ret := _m.Called(ctx, trainerUUID)

	var r0 []trainer.WorkoutGroup
	if rf, ok := ret.Get(0).(func(context.Context, string) []trainer.WorkoutGroup); ok {
		r0 = rf(ctx, trainerUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]trainer.WorkoutGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, trainerUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryWorkoutGroups'
type WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call struct {
	*mock.Call
}

// QueryWorkoutGroups is a helper method to define mock.On call
//  - ctx context.Context
//  - trainerUUID string
func (_e *WorkoutGroupsHandlerRepository_Expecter) QueryWorkoutGroups(ctx interface{}, trainerUUID interface{}) *WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call {
	return &WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call{Call: _e.mock.On("QueryWorkoutGroups", ctx, trainerUUID)}
}

func (_c *WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call) Run(run func(ctx context.Context, trainerUUID string)) *WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call) Return(_a0 []trainer.WorkoutGroup, _a1 error) *WorkoutGroupsHandlerRepository_QueryWorkoutGroups_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}
