package cart

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/dto"
	"github.com/labstack/echo/v4"
)

type RepositoryCartInterface interface {
	CreateCart(newCart *entities.CartModels) (*entities.CartModels, error)
	CreateCartItem(cartItem *entities.CartItemModels) (*entities.CartItemModels, error)
	GetCartItemByProductID(cartID, productID uint64) (*entities.CartItemModels, error)
	GetCartItemsByCartID(cartID uint64) ([]*entities.CartItemModels, error)
	GetCartItemByID(cartItemID uint64) (*entities.CartItemModels, error)
	GetCartByID(cartID uint64) (*entities.CartModels, error)
	GetCartByUserID(userID uint64) (*entities.CartModels, error)
	GetCart(userID uint64) (*entities.CartModels, error)
	UpdateCartItem(cartItem *entities.CartItemModels) error
	UpdateGrandTotal(cartID, grandTotal uint64) error
	DeleteCartItem(cartItemID uint64) error
}

type ServiceCartInterface interface {
	AddCartItems(userID uint64, request *dto.AddCartItemsRequest) (*entities.CartItemModels, error)
	GetCart(userID uint64) (*entities.CartModels, error)
	ReduceCartItemQuantity(cartItemID, quantity uint64) error
	DeleteCartItem(cartItemID uint64) error
}

type HandlerCartInterface interface {
	AddCartItem() echo.HandlerFunc
	GetCart() echo.HandlerFunc
	ReduceQuantity() echo.HandlerFunc
	DeleteCartItems() echo.HandlerFunc
}