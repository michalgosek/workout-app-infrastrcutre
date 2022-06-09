// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	rest "github.com/michalgosek/workout-app-infrastrcutre/service-discovery/internal/rest"
	mock "github.com/stretchr/testify/mock"
)

// RegisterHandlerOption is an autogenerated mock type for the RegisterHandlerOption type
type RegisterHandlerOption struct {
	mock.Mock
}

type RegisterHandlerOption_Expecter struct {
	mock *mock.Mock
}

func (_m *RegisterHandlerOption) EXPECT() *RegisterHandlerOption_Expecter {
	return &RegisterHandlerOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: r
func (_m *RegisterHandlerOption) Execute(r *rest.RegisterHandler) {
	_m.Called(r)
}

// RegisterHandlerOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type RegisterHandlerOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//  - r *rest.RegisterHandler
func (_e *RegisterHandlerOption_Expecter) Execute(r interface{}) *RegisterHandlerOption_Execute_Call {
	return &RegisterHandlerOption_Execute_Call{Call: _e.mock.On("Execute", r)}
}

func (_c *RegisterHandlerOption_Execute_Call) Run(run func(r *rest.RegisterHandler)) *RegisterHandlerOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*rest.RegisterHandler))
	})
	return _c
}

func (_c *RegisterHandlerOption_Execute_Call) Return() *RegisterHandlerOption_Execute_Call {
	_c.Call.Return()
	return _c
}
