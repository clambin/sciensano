// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	bracket "github.com/clambin/sciensano/internal/population/bracket"

	mock "github.com/stretchr/testify/mock"
)

// Fetcher is an autogenerated mock type for the Fetcher type
type Fetcher struct {
	mock.Mock
}

type Fetcher_Expecter struct {
	mock *mock.Mock
}

func (_m *Fetcher) EXPECT() *Fetcher_Expecter {
	return &Fetcher_Expecter{mock: &_m.Mock}
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

// Fetcher_GetByAgeBracket_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByAgeBracket'
type Fetcher_GetByAgeBracket_Call struct {
	*mock.Call
}

// GetByAgeBracket is a helper method to define mock.On call
//   - arguments bracket.Bracket
func (_e *Fetcher_Expecter) GetByAgeBracket(arguments interface{}) *Fetcher_GetByAgeBracket_Call {
	return &Fetcher_GetByAgeBracket_Call{Call: _e.mock.On("GetByAgeBracket", arguments)}
}

func (_c *Fetcher_GetByAgeBracket_Call) Run(run func(arguments bracket.Bracket)) *Fetcher_GetByAgeBracket_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bracket.Bracket))
	})
	return _c
}

func (_c *Fetcher_GetByAgeBracket_Call) Return(count int) *Fetcher_GetByAgeBracket_Call {
	_c.Call.Return(count)
	return _c
}

func (_c *Fetcher_GetByAgeBracket_Call) RunAndReturn(run func(bracket.Bracket) int) *Fetcher_GetByAgeBracket_Call {
	_c.Call.Return(run)
	return _c
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

// Fetcher_GetByRegion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByRegion'
type Fetcher_GetByRegion_Call struct {
	*mock.Call
}

// GetByRegion is a helper method to define mock.On call
func (_e *Fetcher_Expecter) GetByRegion() *Fetcher_GetByRegion_Call {
	return &Fetcher_GetByRegion_Call{Call: _e.mock.On("GetByRegion")}
}

func (_c *Fetcher_GetByRegion_Call) Run(run func()) *Fetcher_GetByRegion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fetcher_GetByRegion_Call) Return(figures map[string]int) *Fetcher_GetByRegion_Call {
	_c.Call.Return(figures)
	return _c
}

func (_c *Fetcher_GetByRegion_Call) RunAndReturn(run func() map[string]int) *Fetcher_GetByRegion_Call {
	_c.Call.Return(run)
	return _c
}

// WaitTillReady provides a mock function with given fields: ctx
func (_m *Fetcher) WaitTillReady(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetcher_WaitTillReady_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WaitTillReady'
type Fetcher_WaitTillReady_Call struct {
	*mock.Call
}

// WaitTillReady is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Fetcher_Expecter) WaitTillReady(ctx interface{}) *Fetcher_WaitTillReady_Call {
	return &Fetcher_WaitTillReady_Call{Call: _e.mock.On("WaitTillReady", ctx)}
}

func (_c *Fetcher_WaitTillReady_Call) Run(run func(ctx context.Context)) *Fetcher_WaitTillReady_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Fetcher_WaitTillReady_Call) Return(_a0 error) *Fetcher_WaitTillReady_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Fetcher_WaitTillReady_Call) RunAndReturn(run func(context.Context) error) *Fetcher_WaitTillReady_Call {
	_c.Call.Return(run)
	return _c
}

// NewFetcher creates a new instance of Fetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *Fetcher {
	mock := &Fetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
