// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	apiclient "github.com/clambin/sciensano/apiclient"

	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Holder is an autogenerated mock type for the Holder type
type Holder struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *Holder) Get(name string) ([]apiclient.APIResponse, bool) {
	ret := _m.Called(name)

	var r0 []apiclient.APIResponse
	if rf, ok := ret.Get(0).(func(string) []apiclient.APIResponse); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]apiclient.APIResponse)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Run provides a mock function with given fields: ctx, interval
func (_m *Holder) Run(ctx context.Context, interval time.Duration) {
	_m.Called(ctx, interval)
}

// Stats provides a mock function with given fields:
func (_m *Holder) Stats() map[string]int {
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
