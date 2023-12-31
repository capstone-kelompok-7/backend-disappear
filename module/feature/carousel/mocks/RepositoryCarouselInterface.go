// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	entities "github.com/capstone-kelompok-7/backend-disappear/module/entities"
	mock "github.com/stretchr/testify/mock"
)

// RepositoryCarouselInterface is an autogenerated mock type for the RepositoryCarouselInterface type
type RepositoryCarouselInterface struct {
	mock.Mock
}

// CreateCarousel provides a mock function with given fields: _a0
func (_m *RepositoryCarouselInterface) CreateCarousel(_a0 *entities.CarouselModels) (*entities.CarouselModels, error) {
	ret := _m.Called(_a0)

	var r0 *entities.CarouselModels
	var r1 error
	if rf, ok := ret.Get(0).(func(*entities.CarouselModels) (*entities.CarouselModels, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*entities.CarouselModels) *entities.CarouselModels); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.CarouselModels)
		}
	}

	if rf, ok := ret.Get(1).(func(*entities.CarouselModels) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCarousel provides a mock function with given fields: carouselID
func (_m *RepositoryCarouselInterface) DeleteCarousel(carouselID uint64) error {
	ret := _m.Called(carouselID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(carouselID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: page, perPage
func (_m *RepositoryCarouselInterface) FindAll(page int, perPage int) ([]*entities.CarouselModels, error) {
	ret := _m.Called(page, perPage)

	var r0 []*entities.CarouselModels
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) ([]*entities.CarouselModels, error)); ok {
		return rf(page, perPage)
	}
	if rf, ok := ret.Get(0).(func(int, int) []*entities.CarouselModels); ok {
		r0 = rf(page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.CarouselModels)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(page, perPage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByName provides a mock function with given fields: page, perPage, name
func (_m *RepositoryCarouselInterface) FindByName(page int, perPage int, name string) ([]*entities.CarouselModels, error) {
	ret := _m.Called(page, perPage, name)

	var r0 []*entities.CarouselModels
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, string) ([]*entities.CarouselModels, error)); ok {
		return rf(page, perPage, name)
	}
	if rf, ok := ret.Get(0).(func(int, int, string) []*entities.CarouselModels); ok {
		r0 = rf(page, perPage, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.CarouselModels)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, string) error); ok {
		r1 = rf(page, perPage, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCarouselById provides a mock function with given fields: id
func (_m *RepositoryCarouselInterface) GetCarouselById(id uint64) (*entities.CarouselModels, error) {
	ret := _m.Called(id)

	var r0 *entities.CarouselModels
	var r1 error
	if rf, ok := ret.Get(0).(func(uint64) (*entities.CarouselModels, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint64) *entities.CarouselModels); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.CarouselModels)
		}
	}

	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalCarouselCount provides a mock function with given fields:
func (_m *RepositoryCarouselInterface) GetTotalCarouselCount() (int64, error) {
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

// GetTotalCarouselCountByName provides a mock function with given fields: name
func (_m *RepositoryCarouselInterface) GetTotalCarouselCountByName(name string) (int64, error) {
	ret := _m.Called(name)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCarousel provides a mock function with given fields: carouselID, updatedCarousel
func (_m *RepositoryCarouselInterface) UpdateCarousel(carouselID uint64, updatedCarousel *entities.CarouselModels) error {
	ret := _m.Called(carouselID, updatedCarousel)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, *entities.CarouselModels) error); ok {
		r0 = rf(carouselID, updatedCarousel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepositoryCarouselInterface creates a new instance of RepositoryCarouselInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepositoryCarouselInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepositoryCarouselInterface {
	mock := &RepositoryCarouselInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
