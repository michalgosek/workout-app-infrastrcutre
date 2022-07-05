// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	trainer "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainer"
	mock "github.com/stretchr/testify/mock"
)

// TrainerRepository is an autogenerated mock type for the TrainerRepository type
type TrainerRepository struct {
	mock.Mock
}

type TrainerRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *TrainerRepository) EXPECT() *TrainerRepository_Expecter {
	return &TrainerRepository_Expecter{mock: &_m.Mock}
}

// QueryTrainerWorkoutGroup provides a mock function with given fields: ctx, trainerUUID, groupUUID
func (_m *TrainerRepository) QueryTrainerWorkoutGroup(ctx context.Context, trainerUUID string, groupUUID string) (trainer.WorkoutGroup, error) {
	ret := _m.Called(ctx, trainerUUID, groupUUID)

	var r0 trainer.WorkoutGroup
	if rf, ok := ret.Get(0).(func(context.Context, string, string) trainer.WorkoutGroup); ok {
		r0 = rf(ctx, trainerUUID, groupUUID)
	} else {
		r0 = ret.Get(0).(trainer.WorkoutGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, trainerUUID, groupUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TrainerRepository_QueryTrainerWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryTrainerWorkoutGroup'
type TrainerRepository_QueryTrainerWorkoutGroup_Call struct {
	*mock.Call
}

// QueryTrainerWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - trainerUUID string
//  - groupUUID string
func (_e *TrainerRepository_Expecter) QueryTrainerWorkoutGroup(ctx interface{}, trainerUUID interface{}, groupUUID interface{}) *TrainerRepository_QueryTrainerWorkoutGroup_Call {
	return &TrainerRepository_QueryTrainerWorkoutGroup_Call{Call: _e.mock.On("QueryTrainerWorkoutGroup", ctx, trainerUUID, groupUUID)}
}

func (_c *TrainerRepository_QueryTrainerWorkoutGroup_Call) Run(run func(ctx context.Context, trainerUUID string, groupUUID string)) *TrainerRepository_QueryTrainerWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TrainerRepository_QueryTrainerWorkoutGroup_Call) Return(_a0 trainer.WorkoutGroup, _a1 error) *TrainerRepository_QueryTrainerWorkoutGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// UpsertTrainerWorkoutGroup provides a mock function with given fields: ctx, group
func (_m *TrainerRepository) UpsertTrainerWorkoutGroup(ctx context.Context, group trainer.WorkoutGroup) error {
	ret := _m.Called(ctx, group)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, trainer.WorkoutGroup) error); ok {
		r0 = rf(ctx, group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TrainerRepository_UpsertTrainerWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertTrainerWorkoutGroup'
type TrainerRepository_UpsertTrainerWorkoutGroup_Call struct {
	*mock.Call
}

// UpsertTrainerWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - group trainer.WorkoutGroup
func (_e *TrainerRepository_Expecter) UpsertTrainerWorkoutGroup(ctx interface{}, group interface{}) *TrainerRepository_UpsertTrainerWorkoutGroup_Call {
	return &TrainerRepository_UpsertTrainerWorkoutGroup_Call{Call: _e.mock.On("UpsertTrainerWorkoutGroup", ctx, group)}
}

func (_c *TrainerRepository_UpsertTrainerWorkoutGroup_Call) Run(run func(ctx context.Context, group trainer.WorkoutGroup)) *TrainerRepository_UpsertTrainerWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(trainer.WorkoutGroup))
	})
	return _c
}

func (_c *TrainerRepository_UpsertTrainerWorkoutGroup_Call) Return(_a0 error) *TrainerRepository_UpsertTrainerWorkoutGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewTrainerRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTrainerRepository creates a new instance of TrainerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTrainerRepository(t mockConstructorTestingTNewTrainerRepository) *TrainerRepository {
	mock := &TrainerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
