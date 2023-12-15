package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
)

type CartService struct {
	repo           cart.RepositoryCartInterface
	productService product.ServiceProductInterface
}

func NewCartService(repo cart.RepositoryCartInterface, productService product.ServiceProductInterface) cart.ServiceCartInterface {
	return &CartService{
		repo:           repo,
		productService: productService,
	}
}

func (s *CartService) GetCart(userID uint64) (*entities.CartModels, error) {
	carts, err := s.repo.GetCart(userID)
	if err != nil {
		return nil, errors.New("keranjang tidak ditemukan")
	}
	return carts, nil
}

func (s *CartService) AddCartItems(userID uint64, request *dto.AddCartItemsRequest) (*entities.CartItemModels, error) {
	carts, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		if carts == nil {
			newCart := &entities.CartModels{
				UserID: userID,
			}
			_, err := s.repo.CreateCart(newCart)
			if err != nil {
				return nil, errors.New("gagal membuat keranjang")
			}
			carts = newCart
		}
	}

	existingCartItem, err := s.repo.GetCartItemByProductID(carts.ID, request.ProductID)
	if err == nil && existingCartItem != nil {
		existingCartItem.Quantity += request.Quantity
		existingCartItem.TotalPrice = existingCartItem.Quantity * existingCartItem.Price

		err := s.repo.UpdateCartItem(existingCartItem)
		if err != nil {
			return nil, errors.New("gagal mengubah jumlah produk di keranjang")
		}
		err = s.RecalculateGrandTotal(carts)
		if err != nil {
			return nil, errors.New("gagal menghitung ulang grand total")
		}
		return existingCartItem, nil
	}

	getProductByID, err := s.productService.GetProductByID(request.ProductID)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	cartItem := &entities.CartItemModels{
		CartID:     carts.ID,
		ProductID:  request.ProductID,
		Quantity:   request.Quantity,
		Price:      getProductByID.Price,
		TotalPrice: getProductByID.Price * request.Quantity,
	}

	result, err := s.repo.CreateCartItem(cartItem)
	if err != nil {
		return nil, errors.New("gagal menambahkan produk ke keranjang")
	}

	err = s.RecalculateGrandTotal(carts)
	if err != nil {
		return nil, errors.New("gagal menghitung ulang grand total")
	}

	return result, nil
}

func (s *CartService) DeleteCartItem(cartItemID uint64) error {
	cartItem, err := s.repo.GetCartItemByID(cartItemID)
	if err != nil {
		return errors.New("item dikeranjang tidak ditemukan")
	}

	carts, err := s.repo.GetCartByID(cartItem.CartID)
	if err != nil {
		return errors.New("keranjang tidak di temukan")
	}

	err = s.repo.DeleteCartItem(cartItem.ID)
	if err != nil {
		return errors.New("gagal menghapus item dikeranjang")
	}

	err = s.RecalculateGrandTotal(carts)
	if err != nil {
		return errors.New("gagal untuk menghitung ulang grand total")
	}

	return nil
}

func (s *CartService) ReduceCartItemQuantity(cartItemID, quantity uint64) error {
	cartItem, err := s.repo.GetCartItemByID(cartItemID)
	if err != nil {
		return errors.New("item dikeranjang tidak ditemukan")
	}

	if quantity > cartItem.Quantity {
		return errors.New("jumlah kuantitas yang diminta melebihi jumlah kuantitas yang ada di keranjang")
	}

	cartItem.Quantity -= quantity
	cartItem.TotalPrice = cartItem.Quantity * cartItem.Price

	if cartItem.Quantity == 0 {
		err := s.repo.DeleteCartItem(cartItemID)
		if err != nil {
			return errors.New("gagal menghapus item dikeranjang")
		}
	} else {
		err = s.repo.UpdateCartItem(cartItem)
		if err != nil {
			return errors.New("gagal memperbarui item dikeranjang")
		}
	}

	carts, err := s.repo.GetCartByID(cartItem.CartID)
	if err != nil {
		return errors.New("keranjang tidak ditemukan")
	}
	err = s.RecalculateGrandTotal(carts)
	if err != nil {
		return errors.New("gagal menghitung ulang grand total")
	}

	return nil
}

func (s *CartService) RecalculateGrandTotal(cart *entities.CartModels) error {
	cartItems, err := s.repo.GetCartItemsByCartID(cart.ID)
	if err != nil {
		return err
	}

	var grandTotal uint64
	for _, item := range cartItems {
		grandTotal += item.TotalPrice
	}

	cart.GrandTotal = grandTotal
	err = s.repo.UpdateGrandTotal(cart.ID, grandTotal)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) IsProductInCart(userID, productID uint64) bool {
	isInCart := s.repo.IsProductInCart(userID, productID)
	return isInCart
}

func (s *CartService) RemoveProductFromCart(userID, productID uint64) error {
	isProductInCart := s.repo.IsProductInCart(userID, productID)
	if !isProductInCart {
		return errors.New("produk tidak ada dalam keranjang pengguna")
	}

	err := s.repo.RemoveProductFromCart(userID, productID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) GetCartItems(cartItem uint64) (*entities.CartItemModels, error) {
	cartItems, err := s.repo.GetCartItemByID(cartItem)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}
