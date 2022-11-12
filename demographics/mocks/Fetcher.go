// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	bracket "github.com/clambin/sciensano/demographics/bracket"

	mock "github.com/stretchr/testify/mock"
)

// Fetcher is an autogenerated mock type for the Fetcher type
type Fetcher struct {
	mock.Mock
}

// GetByAgeBracket provides a mock function with given fields: arguments
func (_m *Fetcher) GetByAgeBracket(arguments bracket.Bracket) int {
	ret := _m.Called(arguments)

	var r0 int
	if rf, ok := ret.Get(0).(func(bracket.Bracket) int); ok {
		r0 = rf(arguments)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetByRegion provides a mock function with given fields:
func (_m *Fetcher) GetByRegion() map[string]int {
	ret := _m.Called()

	var r0 map[string]int
	if rf, ok := ret.Get(0).(func() map[string]int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int)
		}
	}

	return r0
}

// Run provides a mock function with given fields: ctx
func (_m *Fetcher) Run(ctx context.Context) {
	_m.Called(ctx)
}

type mockConstructorTestingTNewFetcher interface {
	mock.TestingT
	Cleanup(func())
}

// NewFetcher creates a new instance of Fetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFetcher(t mockConstructorTestingTNewFetcher) *Fetcher {
	mock := &Fetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
