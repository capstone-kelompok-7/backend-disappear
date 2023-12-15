// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	dto "github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/dto"

	mock "github.com/stretchr/testify/mock"
)

// ServiceCartInterface is an autogenerated mock type for the ServiceCartInterface type
type ServiceCartInterface struct {
	mock.Mock
}

// AddCartItems provides a mock function with given fields: userID, request
func (_m *ServiceCartInterface) AddCartItems(userID uint64, request *dto.AddCartItemsRequest) (*entities.CartItemModels, error) {
	ret := _m.Called(userID, request)

	var r0 *entities.CartItemModels
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64, *dto.AddCartItemsRequest) (*entities.CartItemModels, error)); ok {
		return rf(userID, request)
	}
	if rf, ok := ret.Get(0).(func(uint64, *dto.AddCartItemsRequest) *entities.CartItemModels); ok {
		r0 = rf(userID, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.CartItemModels)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64, *dto.AddCartItemsRequest) error); ok {
		r1 = rf(userID, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCartItem provides a mock function with given fields: cartItemID
func (_m *ServiceCartInterface) DeleteCartItem(cartItemID uint64) error {
	ret := _m.Called(cartItemID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(cartItemID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCart provides a mock function with given fields: userID
func (_m *ServiceCartInterface) GetCart(userID uint64) (*entities.CartModels, error) {
	ret := _m.Called(userID)

	var r0 *entities.CartModels
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*entities.CartModels, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(uint64) *entities.CartModels); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.CartModels)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCartItems provides a mock function with given fields: cartItem
func (_m *ServiceCartInterface) GetCartItems(cartItem uint64) (*entities.CartItemModels, error) {
	ret := _m.Called(cartItem)

	var r0 *entities.CartItemModels
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*entities.CartItemModels, error)); ok {
		return rf(cartItem)
	}
	if rf, ok := ret.Get(0).(func(uint64) *entities.CartItemModels); ok {
		r0 = rf(cartItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.CartItemModels)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(cartItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsProductInCart provides a mock function with given fields: userID, productID
func (_m *ServiceCartInterface) IsProductInCart(userID uint64, productID uint64) bool {
	ret := _m.Called(userID, productID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, uint64) bool); ok {
		r0 = rf(userID, productID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RecalculateGrandTotal provides a mock function with given fields: _a0
func (_m *ServiceCartInterface) RecalculateGrandTotal(_a0 *entities.CartModels) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entities.CartModels) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReduceCartItemQuantity provides a mock function with given fields: cartItemID, quantity
func (_m *ServiceCartInterface) ReduceCartItemQuantity(cartItemID uint64, quantity uint64) error {
	ret := _m.Called(cartItemID, quantity)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, uint64) error); ok {
		r0 = rf(cartItemID, quantity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveProductFromCart provides a mock function with given fields: userID, productID
func (_m *ServiceCartInterface) RemoveProductFromCart(userID uint64, productID uint64) error {
	ret := _m.Called(userID, productID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, uint64) error); ok {
		r0 = rf(userID, productID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServiceCartInterface creates a new instance of ServiceCartInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceCartInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceCartInterface {
	mock := &ServiceCartInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
