// Code generated by mockery. DO NOT EDIT.

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.
package injector

import mock "github.com/stretchr/testify/mock"

// MockBPFDiskFailureCommand is an autogenerated mock type for the BPFDiskFailureCommand type
type MockBPFDiskFailureCommand struct {
	mock.Mock
}

type MockBPFDiskFailureCommand_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBPFDiskFailureCommand) EXPECT() *MockBPFDiskFailureCommand_Expecter {
	return &MockBPFDiskFailureCommand_Expecter{mock: &_m.Mock}
}

// Run provides a mock function with given fields: pid, path
func (_m *MockBPFDiskFailureCommand) Run(pid int, path string) error {
	ret := _m.Called(pid, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, string) error); ok {
		r0 = rf(pid, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBPFDiskFailureCommand_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type MockBPFDiskFailureCommand_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - pid int
//   - path string
func (_e *MockBPFDiskFailureCommand_Expecter) Run(pid interface{}, path interface{}) *MockBPFDiskFailureCommand_Run_Call {
	return &MockBPFDiskFailureCommand_Run_Call{Call: _e.mock.On("Run", pid, path)}
}

func (_c *MockBPFDiskFailureCommand_Run_Call) Run(run func(pid int, path string)) *MockBPFDiskFailureCommand_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(string))
	})
	return _c
}

func (_c *MockBPFDiskFailureCommand_Run_Call) Return(_a0 error) *MockBPFDiskFailureCommand_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBPFDiskFailureCommand_Run_Call) RunAndReturn(run func(int, string) error) *MockBPFDiskFailureCommand_Run_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockBPFDiskFailureCommand interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBPFDiskFailureCommand creates a new instance of MockBPFDiskFailureCommand. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBPFDiskFailureCommand(t mockConstructorTestingTNewMockBPFDiskFailureCommand) *MockBPFDiskFailureCommand {
	mock := &MockBPFDiskFailureCommand{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}