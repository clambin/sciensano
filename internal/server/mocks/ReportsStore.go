// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	tabulator "github.com/clambin/go-common/tabulator"
)

// ReportsStore is an autogenerated mock type for the ReportsStore type
type ReportsStore struct {
	mock.Mock
}

type ReportsStore_Expecter struct {
	mock *mock.Mock
}

func (_m *ReportsStore) EXPECT() *ReportsStore_Expecter {
	return &ReportsStore_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: _a0
func (_m *ReportsStore) Get(_a0 string) (*tabulator.Tabulator, error) {
	ret := _m.Called(_a0)

	var r0 *tabulator.Tabulator
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*tabulator.Tabulator, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *tabulator.Tabulator); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tabulator.Tabulator)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReportsStore_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ReportsStore_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 string
func (_e *ReportsStore_Expecter) Get(_a0 interface{}) *ReportsStore_Get_Call {
	return &ReportsStore_Get_Call{Call: _e.mock.On("Get", _a0)}
}

func (_c *ReportsStore_Get_Call) Run(run func(_a0 string)) *ReportsStore_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ReportsStore_Get_Call) Return(_a0 *tabulator.Tabulator, _a1 error) *ReportsStore_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ReportsStore_Get_Call) RunAndReturn(run func(string) (*tabulator.Tabulator, error)) *ReportsStore_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Keys provides a mock function with given fields:
func (_m *ReportsStore) Keys() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// ReportsStore_Keys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Keys'
type ReportsStore_Keys_Call struct {
	*mock.Call
}

// Keys is a helper method to define mock.On call
func (_e *ReportsStore_Expecter) Keys() *ReportsStore_Keys_Call {
	return &ReportsStore_Keys_Call{Call: _e.mock.On("Keys")}
}

func (_c *ReportsStore_Keys_Call) Run(run func()) *ReportsStore_Keys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReportsStore_Keys_Call) Return(_a0 []string) *ReportsStore_Keys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReportsStore_Keys_Call) RunAndReturn(run func() []string) *ReportsStore_Keys_Call {
	_c.Call.Return(run)
	return _c
}

// NewReportsStore creates a new instance of ReportsStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReportsStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReportsStore {
	mock := &ReportsStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
