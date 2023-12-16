package service

import (
	"errors"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/config"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	assistantMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/mocks"
	assistants "github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/service"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/cart/mocks"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/product"
	productsMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/mocks"
	products "github.com/capstone-kelompok-7/backend-disappear/module/feature/product/service"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) (*mocks.RepositoryCartInterface, cart.ServiceCartInterface, product.ServiceProductInterface, assistant.ServiceAssistantInterface) {
	repo := mocks.NewRepositoryCartInterface(t)
	repoProduct := productsMocks.NewRepositoryProductInterface(t)
	repoAssistant := assistantMocks.NewRepositoryAssistantInterface(t)

	assistantService := assistants.NewAssistantService(repoAssistant, nil, config.Config{})
	productService := products.NewProductService(repoProduct, assistantService)
	cartService := NewCartService(repo, productService)

	return repo, cartService, productService, assistantService
}

func createExpectedCart() *entities.CartModels {
	expectedCartItem := &entities.CartItemModels{
		ID:         1,
		CartID:     1,
		ProductID:  1,
		Quantity:   2,
		Price:      5000,
		TotalPrice: 10000,
		Product: &entities.ProductModels{
			ID:            1,
			Name:          "Product A",
			Description:   "Description",
			GramPlastic:   100,
			Price:         5000,
			Stock:         10,
			Discount:      0,
			Exp:           0,
			Rating:        0,
			TotalReview:   0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     nil,
			ProductPhotos: nil,
			ProductReview: nil,
			Categories:    nil,
		},
	}

	expectedCart := &entities.CartModels{
		ID:         1,
		UserID:     1,
		GrandTotal: 10000,
		CartItems:  []*entities.CartItemModels{expectedCartItem}, // Gunakan slice untuk CartItems
	}
	return expectedCart
}

func TestCartService_GetCart(t *testing.T) {
	t.Run("Success Case - Cart Found", func(t *testing.T) {
		userID := uint64(1)
		expectedCart := createExpectedCart()

		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCart", userID).Return(expectedCart, nil)

		carts, err := cartService.GetCart(userID)

		assert.NoError(t, err)
		assert.NotNil(t, carts)
		assert.Equal(t, expectedCart, carts)
	})

	t.Run("Failed Case - Cart Not Found", func(t *testing.T) {
		userID := uint64(1)

		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCart", userID).Return(nil, errors.New("keranjang tidak ditemukan"))

		carts, err := cartService.GetCart(userID)

		assert.Error(t, err)
		assert.Nil(t, carts)
	})
}

func TestCartService_RemoveProductFromCart(t *testing.T) {
	userID := uint64(1)
	productID := uint64(10)

	t.Run("Success Case - Product Removed From Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("IsProductInCart", userID, productID).Return(true)
		repoMock.On("RemoveProductFromCart", userID, productID).Return(nil)

		err := cartService.RemoveProductFromCart(userID, productID)

		assert.NoError(t, err)
	})

	t.Run("Failed Case - Product Not in Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("IsProductInCart", userID, productID).Return(false)

		err := cartService.RemoveProductFromCart(userID, productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "produk tidak ada dalam keranjang pengguna")
	})

	t.Run("Failed Case - Error Removed Product From Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("IsProductInCart", userID, productID).Return(true)
		repoMock.On("RemoveProductFromCart", userID, productID).Return(errors.New("gagal menghapus produk dari keranjang"))

		err := cartService.RemoveProductFromCart(userID, productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal menghapus produk dari keranjang")
	})
}

func TestCartService_IsProductInCart(t *testing.T) {
	userID := uint64(1)
	productID := uint64(10)

	t.Run("Success Case - Product In Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("IsProductInCart", userID, productID).Return(true)

		isInCart := cartService.IsProductInCart(userID, productID)

		assert.True(t, isInCart)
	})

	t.Run("Failed Case - Product Not Found in Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("IsProductInCart", userID, productID).Return(false)

		isInCart := cartService.IsProductInCart(userID, productID)

		assert.False(t, isInCart)
	})
}

func TestCartService_GetCartItems(t *testing.T) {
	cartItemID := uint64(1)

	t.Run("Success Case - Cart Items Found ", func(t *testing.T) {
		expectedCartItem := &entities.CartItemModels{
			ID:         cartItemID,
			CartID:     1,
			ProductID:  1,
			Quantity:   2,
			Price:      5000,
			TotalPrice: 10000,
		}

		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCartItem, nil)

		cartItem, err := cartService.GetCartItems(cartItemID)

		assert.NoError(t, err)
		assert.NotNil(t, cartItem)
		assert.Equal(t, expectedCartItem, cartItem)
	})

	t.Run("Failed Case - Cart Item Not Found", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(nil, errors.New("cart item not found"))

		cartItem, err := cartService.GetCartItems(cartItemID)

		assert.Error(t, err)
		assert.Nil(t, cartItem)
		assert.EqualError(t, err, "cart item not found")
	})
}

func TestCartService_RecalculateGrandTotal(t *testing.T) {
	cartID := uint64(1)
	expectedCart := createExpectedCart()

	t.Run("Success Case - Grand Total Recalculated", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemsByCartID", cartID).Return(expectedCart.CartItems, nil)
		repoMock.On("UpdateGrandTotal", cartID, expectedCart.GrandTotal).Return(nil)

		err := cartService.RecalculateGrandTotal(expectedCart)

		assert.NoError(t, err)
		assert.Equal(t, expectedCart.GrandTotal, uint64(10000))
	})

	t.Run("Failed Case - Error Getting Cart Items", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemsByCartID", cartID).Return(nil, errors.New("failed to fetch cart items"))

		err := cartService.RecalculateGrandTotal(expectedCart)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to fetch cart items")
	})

	t.Run("Failed Case - Error Updating Grand Total", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemsByCartID", cartID).Return(expectedCart.CartItems, nil)
		repoMock.On("UpdateGrandTotal", cartID, expectedCart.GrandTotal).Return(errors.New("failed to update grand total"))

		err := cartService.RecalculateGrandTotal(expectedCart)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed to update grand total")
	})
}

func TestCartService_DeleteCartItem(t *testing.T) {
	cartItemID := uint64(1)
	expectedCart := createExpectedCart()

	t.Run("Success Case - Cart Item Deleted", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		cartItem := expectedCart.CartItems[0]

		repoMock.On("GetCartItemByID", cartItem.ID).Return(cartItem, nil).Once()
		repoMock.On("GetCartByID", cartItem.CartID).Return(expectedCart, nil).Once()
		repoMock.On("DeleteCartItem", cartItem.ID).Return(nil).Once()
		repoMock.On("GetCartItemsByCartID", expectedCart.ID).Return(expectedCart.CartItems, nil).Once()
		repoMock.On("UpdateGrandTotal", expectedCart.ID, expectedCart.GrandTotal).Return(nil).Once()

		err := cartService.DeleteCartItem(cartItem.ID)

		assert.NoError(t, err)
	})

	t.Run("Failed Case - Error Getting Cart Item", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(nil, errors.New("failed to fetch cart item"))

		err := cartService.DeleteCartItem(cartItemID)

		assert.Error(t, err)
		assert.EqualError(t, err, "item dikeranjang tidak ditemukan")
	})

	t.Run("Failed Case - Error Getting Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCart.CartItems[0], nil)
		repoMock.On("GetCartByID", expectedCart.ID).Return(nil, errors.New("failed to fetch cart"))

		err := cartService.DeleteCartItem(cartItemID)

		assert.Error(t, err)
		assert.EqualError(t, err, "keranjang tidak di temukan")
	})

	t.Run("Failed Case - Error Deleting Cart Item", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCart.CartItems[0], nil)
		repoMock.On("GetCartByID", expectedCart.ID).Return(expectedCart, nil)
		repoMock.On("DeleteCartItem", cartItemID).Return(errors.New("failed to delete cart item"))

		err := cartService.DeleteCartItem(cartItemID)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal menghapus item dikeranjang")
	})

	t.Run("Failed Case - Error Recalculating Grand Total", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCart.CartItems[0], nil)
		repoMock.On("GetCartByID", expectedCart.ID).Return(expectedCart, nil)
		repoMock.On("DeleteCartItem", cartItemID).Return(nil)
		repoMock.On("GetCartItemsByCartID", expectedCart.ID).Return(expectedCart.CartItems, nil).Once()
		repoMock.On("UpdateGrandTotal", expectedCart.ID, expectedCart.GrandTotal).Return(errors.New("gagal untuk menghitung ulang grand total")).Once()

		err := cartService.DeleteCartItem(cartItemID)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal untuk menghitung ulang grand total")
	})
}

func TestCartService_ReduceCartItemQuantity(t *testing.T) {
	quantity := uint64(2)
	cartItemID := uint64(1)
	expectedCart := createExpectedCart()
	expectedCartItem := &entities.CartItemModels{
		ID:         cartItemID,
		CartID:     1,
		ProductID:  1,
		Quantity:   2,
		Price:      5000,
		TotalPrice: 10000,
	}

	t.Run("Success Case - Reduce Cart Item Quantity", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		expectedCartItem.Quantity = 2
		repoMock.On("GetCartByID", expectedCartItem.CartID).Return(expectedCart, nil)
		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCartItem, nil)
		repoMock.On("GetCartItemsByCartID", expectedCart.ID).Return(expectedCart.CartItems, nil).Once()
		repoMock.On("UpdateGrandTotal", expectedCart.ID, expectedCart.GrandTotal).Return(nil).Once()
		repoMock.On("DeleteCartItem", cartItemID).Return(nil)

		err := cartService.ReduceCartItemQuantity(cartItemID, quantity)

		assert.NoError(t, err)
	})

	t.Run("Failed Case - Item Not Found in Cart", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(nil, errors.New("item dikeranjang tidak ditemukan"))

		err := cartService.ReduceCartItemQuantity(cartItemID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, "item dikeranjang tidak ditemukan")
	})

	t.Run("Failed Case - Requested Quantity Exceeds Cart Item Quantity", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCartItem, nil)

		err := cartService.ReduceCartItemQuantity(cartItemID, expectedCartItem.Quantity+1)

		assert.Error(t, err)
		assert.EqualError(t, err, "jumlah kuantitas yang diminta melebihi jumlah kuantitas yang ada di keranjang")
	})

	t.Run("Failed Case - Error Updating Cart Item", func(t *testing.T) {
		repoMock, cartService, _, _ := setupTestService(t)
		defer repoMock.AssertExpectations(t)

		expectedCartItem.Quantity = 3
		repoMock.On("GetCartItemByID", cartItemID).Return(expectedCartItem, nil)
		repoMock.On("UpdateCartItem", expectedCartItem).Return(errors.New("gagal memperbarui item dikeranjang"))

		err := cartService.ReduceCartItemQuantity(cartItemID, quantity)

		assert.Error(t, err)
		assert.EqualError(t, err, "gagal memperbarui item dikeranjang")
	})

}
