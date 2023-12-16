// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryAssistantInterface is an autogenerated mock type for the RepositoryAssistantInterface type
type RepositoryAssistantInterface struct {
	mock.Mock
}

// CreateAnswer provides a mock function with given fields: chat
func (_m *RepositoryAssistantInterface) CreateAnswer(chat entities.ChatModel) error {
	ret := _m.Called(chat)

	var r0 error
	if rf, ok := ret.Get(0).(func(entities.ChatModel) error); ok {
		r0 = rf(chat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateQuestion provides a mock function with given fields: chat
func (_m *RepositoryAssistantInterface) CreateQuestion(chat entities.ChatModel) error {
	ret := _m.Called(chat)

	var r0 error
	if rf, ok := ret.Get(0).(func(entities.ChatModel) error); ok {
		r0 = rf(chat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetChatByIdUser provides a mock function with given fields: id
func (_m *RepositoryAssistantInterface) GetChatByIdUser(id uint64) ([]entities.ChatModel, error) {
	ret := _m.Called(id)

	var r0 []entities.ChatModel
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) ([]entities.ChatModel, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) []entities.ChatModel); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.ChatModel)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastOrdersByUserID provides a mock function with given fields: userID
func (_m *RepositoryAssistantInterface) GetLastOrdersByUserID(userID uint64) ([]*entities.OrderModels, error) {
	ret := _m.Called(userID)

	var r0 []*entities.OrderModels
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) ([]*entities.OrderModels, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uint64) []*entities.OrderModels); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.OrderModels)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTopRatedProducts provides a mock function with given fields:
func (_m *RepositoryAssistantInterface) GetTopRatedProducts() ([]*entities.ProductModels, error) {
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

// GetTopSellingProducts provides a mock function with given fields:
func (_m *RepositoryAssistantInterface) GetTopSellingProducts() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepositoryAssistantInterface creates a new instance of RepositoryAssistantInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepositoryAssistantInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepositoryAssistantInterface {
	mock := &RepositoryAssistantInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
