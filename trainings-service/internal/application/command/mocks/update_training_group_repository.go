// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	context "context"

	trainings "github.com/michalgosek/workout-app-infrastrcutre/trainings-service/internal/domain/trainings"
	mock "github.com/stretchr/testify/mock"
)

// UpdateTrainingGroupRepository is an autogenerated mock type for the UpdateTrainingGroupRepository type
type UpdateTrainingGroupRepository struct {
	mock.Mock
}

type UpdateTrainingGroupRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UpdateTrainingGroupRepository) EXPECT() *UpdateTrainingGroupRepository_Expecter {
	return &UpdateTrainingGroupRepository_Expecter{mock: &_m.Mock}
}

// UpdateTrainingGroup provides a mock function with given fields: ctx, g
func (_m *UpdateTrainingGroupRepository) UpdateTrainingGroup(ctx context.Context, g *trainings.TrainingGroup) error {
	ret := _m.Called(ctx, g)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *trainings.TrainingGroup) error); ok {
		r0 = rf(ctx, g)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTrainingGroupRepository_UpdateTrainingGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateTrainingGroup'
type UpdateTrainingGroupRepository_UpdateTrainingGroup_Call struct {
	*mock.Call
}

// UpdateTrainingGroup is a helper method to define mock.On call
//  - ctx context.Context
//  - g *trainings.TrainingGroup
func (_e *UpdateTrainingGroupRepository_Expecter) UpdateTrainingGroup(ctx interface{}, g interface{}) *UpdateTrainingGroupRepository_UpdateTrainingGroup_Call {
	return &UpdateTrainingGroupRepository_UpdateTrainingGroup_Call{Call: _e.mock.On("UpdateTrainingGroup", ctx, g)}
}

func (_c *UpdateTrainingGroupRepository_UpdateTrainingGroup_Call) Run(run func(ctx context.Context, g *trainings.TrainingGroup)) *UpdateTrainingGroupRepository_UpdateTrainingGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*trainings.TrainingGroup))
	})
	return _c
}

func (_c *UpdateTrainingGroupRepository_UpdateTrainingGroup_Call) Return(_a0 error) *UpdateTrainingGroupRepository_UpdateTrainingGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewUpdateTrainingGroupRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpdateTrainingGroupRepository creates a new instance of UpdateTrainingGroupRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpdateTrainingGroupRepository(t mockConstructorTestingTNewUpdateTrainingGroupRepository) *UpdateTrainingGroupRepository {
	mock := &UpdateTrainingGroupRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
