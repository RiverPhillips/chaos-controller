// Code generated by mockery. DO NOT EDIT.

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.
package ddmark

import (
	multierror "github.com/hashicorp/go-multierror"
	mock "github.com/stretchr/testify/mock"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

type MockClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClient) EXPECT() *MockClient_Expecter {
	return &MockClient_Expecter{mock: &_m.Mock}
}

// CleanupLibraries provides a mock function with given fields:
func (_m *MockClient) CleanupLibraries() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockClient_CleanupLibraries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CleanupLibraries'
type MockClient_CleanupLibraries_Call struct {
	*mock.Call
}

// CleanupLibraries is a helper method to define mock.On call
func (_e *MockClient_Expecter) CleanupLibraries() *MockClient_CleanupLibraries_Call {
	return &MockClient_CleanupLibraries_Call{Call: _e.mock.On("CleanupLibraries")}
}

func (_c *MockClient_CleanupLibraries_Call) Run(run func()) *MockClient_CleanupLibraries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClient_CleanupLibraries_Call) Return(_a0 error) *MockClient_CleanupLibraries_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClient_CleanupLibraries_Call) RunAndReturn(run func() error) *MockClient_CleanupLibraries_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateStruct provides a mock function with given fields: marshalledStruct, filePath
func (_m *MockClient) ValidateStruct(marshalledStruct interface{}, filePath string) []error {
	ret := _m.Called(marshalledStruct, filePath)

	var r0 []error
	if rf, ok := ret.Get(0).(func(interface{}, string) []error); ok {
		r0 = rf(marshalledStruct, filePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]error)
		}
	}

	return r0
}

// MockClient_ValidateStruct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateStruct'
type MockClient_ValidateStruct_Call struct {
	*mock.Call
}

// ValidateStruct is a helper method to define mock.On call
//   - marshalledStruct interface{}
//   - filePath string
func (_e *MockClient_Expecter) ValidateStruct(marshalledStruct interface{}, filePath interface{}) *MockClient_ValidateStruct_Call {
	return &MockClient_ValidateStruct_Call{Call: _e.mock.On("ValidateStruct", marshalledStruct, filePath)}
}

func (_c *MockClient_ValidateStruct_Call) Run(run func(marshalledStruct interface{}, filePath string)) *MockClient_ValidateStruct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(string))
	})
	return _c
}

func (_c *MockClient_ValidateStruct_Call) Return(_a0 []error) *MockClient_ValidateStruct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClient_ValidateStruct_Call) RunAndReturn(run func(interface{}, string) []error) *MockClient_ValidateStruct_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateStructMultierror provides a mock function with given fields: marshalledStruct, filePath
func (_m *MockClient) ValidateStructMultierror(marshalledStruct interface{}, filePath string) *multierror.Error {
	ret := _m.Called(marshalledStruct, filePath)

	var r0 *multierror.Error
	if rf, ok := ret.Get(0).(func(interface{}, string) *multierror.Error); ok {
		r0 = rf(marshalledStruct, filePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*multierror.Error)
		}
	}

	return r0
}

// MockClient_ValidateStructMultierror_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateStructMultierror'
type MockClient_ValidateStructMultierror_Call struct {
	*mock.Call
}

// ValidateStructMultierror is a helper method to define mock.On call
//   - marshalledStruct interface{}
//   - filePath string
func (_e *MockClient_Expecter) ValidateStructMultierror(marshalledStruct interface{}, filePath interface{}) *MockClient_ValidateStructMultierror_Call {
	return &MockClient_ValidateStructMultierror_Call{Call: _e.mock.On("ValidateStructMultierror", marshalledStruct, filePath)}
}

func (_c *MockClient_ValidateStructMultierror_Call) Run(run func(marshalledStruct interface{}, filePath string)) *MockClient_ValidateStructMultierror_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(string))
	})
	return _c
}

func (_c *MockClient_ValidateStructMultierror_Call) Return(retErr *multierror.Error) *MockClient_ValidateStructMultierror_Call {
	_c.Call.Return(retErr)
	return _c
}

func (_c *MockClient_ValidateStructMultierror_Call) RunAndReturn(run func(interface{}, string) *multierror.Error) *MockClient_ValidateStructMultierror_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockClient(t mockConstructorTestingTNewMockClient) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}