// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	apiclient "github.com/clambin/sciensano/apiclient"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Fetcher is an autogenerated mock type for the Fetcher type
type Fetcher struct {
	mock.Mock
}

// DataTypes provides a mock function with given fields:
func (_m *Fetcher) DataTypes() map[int]string {
	ret := _m.Called()

	var r0 map[int]string
	if rf, ok := ret.Get(0).(func() map[int]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[int]string)
		}
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx, dataType
func (_m *Fetcher) Fetch(ctx context.Context, dataType int) ([]apiclient.APIResponse, error) {
	ret := _m.Called(ctx, dataType)

	var r0 []apiclient.APIResponse
	if rf, ok := ret.Get(0).(func(context.Context, int) []apiclient.APIResponse); ok {
		r0 = rf(ctx, dataType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]apiclient.APIResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, dataType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastUpdated provides a mock function with given fields: ctx, dataType
func (_m *Fetcher) GetLastUpdated(ctx context.Context, dataType int) (time.Time, error) {
	ret := _m.Called(ctx, dataType)

	var r0 time.Time
	if rf, ok := ret.Get(0).(func(context.Context, int) time.Time); ok {
		r0 = rf(ctx, dataType)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, dataType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
