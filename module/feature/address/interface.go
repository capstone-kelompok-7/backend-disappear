package address

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/address/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryAddressInterface interface {
	CreateAddress(newData *entities.AddressModels) (*entities.AddressModels, error)
	FindAllByUserID(userID uint64, page, perPage int) ([]*entities.AddressModels, error)
	GetTotalAddressCountByUserID(userID uint64) (int64, error)
	GetAddressByID(addressID uint64) (*entities.AddressModels, error)
	UpdateAddress(addressID uint64, updatedAddress *dto.UpdateAddressRequest) error
	GetPrimaryAddressByUserID(userID uint64) (*entities.AddressModels, error)
	UpdateIsPrimary(addressID uint64, isPrimary bool) error
	DeleteAddress(addressID uint64) error
}

type ServiceAddressInterface interface {
	CreateAddress(newData *entities.AddressModels) (*entities.AddressModels, error)
	GetAll(userID uint64, page, perPage int) ([]*entities.AddressModels, int64, error)
	CalculatePaginationValues(page int, totalItems int, perPage int) (int, int)
	GetNextPage(currentPage, totalPages int) int
	GetPrevPage(currentPage int) int
	GetAddressByID(addressID uint64) (*entities.AddressModels, error)
	UpdateAddress(userID, addressID uint64, updatedAddress *dto.UpdateAddressRequest) error
	DeleteAddress(addressID, userID uint64) error
}

type HandlerAddressInterface interface {
	CreateAddress() echo.HandlerFunc
	GetAllAddress() echo.HandlerFunc
	UpdateAddress() echo.HandlerFunc
	DeleteAddress() echo.HandlerFunc
}
