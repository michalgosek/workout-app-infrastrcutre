// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	query "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/application/trainer/query"
	mock "github.com/stretchr/testify/mock"
)

// TrainerWorkoutGroupReadModel is an autogenerated mock type for the TrainerWorkoutGroupReadModel type
type TrainerWorkoutGroupReadModel struct {
	mock.Mock
}

type TrainerWorkoutGroupReadModel_Expecter struct {
	mock *mock.Mock
}

func (_m *TrainerWorkoutGroupReadModel) EXPECT() *TrainerWorkoutGroupReadModel_Expecter {
	return &TrainerWorkoutGroupReadModel_Expecter{mock: &_m.Mock}
}

// TrainerWorkoutGroup provides a mock function with given fields: ctx, trainerUUID, groupUUID
func (_m *TrainerWorkoutGroupReadModel) TrainerWorkoutGroup(ctx context.Context, trainerUUID string, groupUUID string) (query.TrainerWorkoutGroup, error) {
	ret := _m.Called(ctx, trainerUUID, groupUUID)

	var r0 query.TrainerWorkoutGroup
	if rf, ok := ret.Get(0).(func(context.Context, string, string) query.TrainerWorkoutGroup); ok {
		r0 = rf(ctx, trainerUUID, groupUUID)
	} else {
		r0 = ret.Get(0).(query.TrainerWorkoutGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, trainerUUID, groupUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TrainerWorkoutGroup'
type TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call struct {
	*mock.Call
}

// TrainerWorkoutGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - trainerUUID string
//  - groupUUID string
func (_e *TrainerWorkoutGroupReadModel_Expecter) TrainerWorkoutGroup(ctx interface{}, trainerUUID interface{}, groupUUID interface{}) *TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call {
	return &TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call{Call: _e.mock.On("TrainerWorkoutGroup", ctx, trainerUUID, groupUUID)}
}

func (_c *TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call) Run(run func(ctx context.Context, trainerUUID string, groupUUID string)) *TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call) Return(_a0 query.TrainerWorkoutGroup, _a1 error) *TrainerWorkoutGroupReadModel_TrainerWorkoutGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewTrainerWorkoutGroupReadModel interface {
	mock.TestingT
	Cleanup(func())
}

// NewTrainerWorkoutGroupReadModel creates a new instance of TrainerWorkoutGroupReadModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTrainerWorkoutGroupReadModel(t mockConstructorTestingTNewTrainerWorkoutGroupReadModel) *TrainerWorkoutGroupReadModel {
	mock := &TrainerWorkoutGroupReadModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
