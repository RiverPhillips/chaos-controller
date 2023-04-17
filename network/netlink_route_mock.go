// Code generated by mockery. DO NOT EDIT.

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.
package network

import (
	net "net"

	mock "github.com/stretchr/testify/mock"
)

// MockNetlinkRoute is an autogenerated mock type for the NetlinkRoute type
type MockNetlinkRoute struct {
	mock.Mock
}

type MockNetlinkRoute_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNetlinkRoute) EXPECT() *MockNetlinkRoute_Expecter {
	return &MockNetlinkRoute_Expecter{mock: &_m.Mock}
}

// Gateway provides a mock function with given fields:
func (_m *MockNetlinkRoute) Gateway() net.IP {
	ret := _m.Called()

	var r0 net.IP
	if rf, ok := ret.Get(0).(func() net.IP); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.IP)
		}
	}

	return r0
}

// MockNetlinkRoute_Gateway_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Gateway'
type MockNetlinkRoute_Gateway_Call struct {
	*mock.Call
}

// Gateway is a helper method to define mock.On call
func (_e *MockNetlinkRoute_Expecter) Gateway() *MockNetlinkRoute_Gateway_Call {
	return &MockNetlinkRoute_Gateway_Call{Call: _e.mock.On("Gateway")}
}

func (_c *MockNetlinkRoute_Gateway_Call) Run(run func()) *MockNetlinkRoute_Gateway_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNetlinkRoute_Gateway_Call) Return(_a0 net.IP) *MockNetlinkRoute_Gateway_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetlinkRoute_Gateway_Call) RunAndReturn(run func() net.IP) *MockNetlinkRoute_Gateway_Call {
	_c.Call.Return(run)
	return _c
}

// Link provides a mock function with given fields:
func (_m *MockNetlinkRoute) Link() NetlinkLink {
	ret := _m.Called()

	var r0 NetlinkLink
	if rf, ok := ret.Get(0).(func() NetlinkLink); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(NetlinkLink)
		}
	}

	return r0
}

// MockNetlinkRoute_Link_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Link'
type MockNetlinkRoute_Link_Call struct {
	*mock.Call
}

// Link is a helper method to define mock.On call
func (_e *MockNetlinkRoute_Expecter) Link() *MockNetlinkRoute_Link_Call {
	return &MockNetlinkRoute_Link_Call{Call: _e.mock.On("Link")}
}

func (_c *MockNetlinkRoute_Link_Call) Run(run func()) *MockNetlinkRoute_Link_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockNetlinkRoute_Link_Call) Return(_a0 NetlinkLink) *MockNetlinkRoute_Link_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetlinkRoute_Link_Call) RunAndReturn(run func() NetlinkLink) *MockNetlinkRoute_Link_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockNetlinkRoute interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockNetlinkRoute creates a new instance of MockNetlinkRoute. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockNetlinkRoute(t mockConstructorTestingTNewMockNetlinkRoute) *MockNetlinkRoute {
	mock := &MockNetlinkRoute{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}