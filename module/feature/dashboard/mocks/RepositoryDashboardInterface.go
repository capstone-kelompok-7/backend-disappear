// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// RepositoryDashboardInterface is an autogenerated mock type for the RepositoryDashboardInterface type
type RepositoryDashboardInterface struct {
	mock.Mock
}

// CountIncome provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) CountIncome() (float64, error) {
	ret := _m.Called()

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func() (float64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountOrder provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) CountOrder() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountProducts provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) CountProducts() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountTotalGram provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) CountTotalGram() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountUsers provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) CountUsers() (int64, error) {
	ret := _m.Called()

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGramPlasticStat provides a mock function with given fields: startOfWeek, endOfWeek
func (_m *RepositoryDashboardInterface) GetGramPlasticStat(startOfWeek time.Time, endOfWeek time.Time) (uint64, error) {
	ret := _m.Called(startOfWeek, endOfWeek)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) (uint64, error)); ok {
		return rf(startOfWeek, endOfWeek)
	}
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) uint64); ok {
		r0 = rf(startOfWeek, endOfWeek)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(startOfWeek, endOfWeek)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestTransactions provides a mock function with given fields: limit
func (_m *RepositoryDashboardInterface) GetLatestTransactions(limit int) ([]*entities.OrderModels, error) {
	ret := _m.Called(limit)

	var r0 []*entities.OrderModels
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]*entities.OrderModels, error)); ok {
		return rf(limit)
	}
	if rf, ok := ret.Get(0).(func(int) []*entities.OrderModels); ok {
		r0 = rf(limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.OrderModels)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductWithMaxReviews provides a mock function with given fields:
func (_m *RepositoryDashboardInterface) GetProductWithMaxReviews() ([]*entities.ProductModels, error) {
	ret := _m.Called()

	var r0 []*entities.ProductModels
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*entities.ProductModels, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*entities.ProductModels); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.ProductModels)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepositoryDashboardInterface creates a new instance of RepositoryDashboardInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepositoryDashboardInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepositoryDashboardInterface {
	mock := &RepositoryDashboardInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
