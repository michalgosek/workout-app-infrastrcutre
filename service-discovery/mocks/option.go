// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	registry "service-discovery/internal/registry"

	mock "github.com/stretchr/testify/mock"
)

// Option is an autogenerated mock type for the Option type
type Option struct {
	mock.Mock
}

type Option_Expecter struct {
	mock *mock.Mock
}

func (_m *Option) EXPECT() *Option_Expecter {
	return &Option_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: s
func (_m *Option) Execute(s *registry.Service) {
	_m.Called(s)
}

// Option_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type Option_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//  - s *registry.Service
func (_e *Option_Expecter) Execute(s interface{}) *Option_Execute_Call {
	return &Option_Execute_Call{Call: _e.mock.On("Execute", s)}
}

func (_c *Option_Execute_Call) Run(run func(s *registry.Service)) *Option_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*registry.Service))
	})
	return _c
}

func (_c *Option_Execute_Call) Return() *Option_Execute_Call {
	_c.Call.Return()
	return _c
}