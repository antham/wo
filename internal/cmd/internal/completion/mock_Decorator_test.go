// Code generated by mockery v2.44.2. DO NOT EDIT.

package completion

import (
	cobra "github.com/spf13/cobra"
	mock "github.com/stretchr/testify/mock"
)

// MockDecorator is an autogenerated mock type for the Decorator type
type MockDecorator struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDecorator) Execute(_a0 workspaceManager, _a1 string, _a2 ...string) ([]string, cobra.ShellCompDirective, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []string
	var r1 cobra.ShellCompDirective
	var r2 error
	if rf, ok := ret.Get(0).(func(workspaceManager, string, ...string) ([]string, cobra.ShellCompDirective, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(workspaceManager, string, ...string) []string); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(workspaceManager, string, ...string) cobra.ShellCompDirective); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Get(1).(cobra.ShellCompDirective)
	}

	if rf, ok := ret.Get(2).(func(workspaceManager, string, ...string) error); ok {
		r2 = rf(_a0, _a1, _a2...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewMockDecorator creates a new instance of MockDecorator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDecorator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDecorator {
	mock := &MockDecorator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
