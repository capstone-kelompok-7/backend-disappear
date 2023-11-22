package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type CartItem struct {
	CartItemID  uint64 `json:"cart_item_id"`
	ProductName string `json:"product_name"`
	GramPlastic uint64 `json:"gram_plastic"`
	Price       uint64 `json:"price"`
	Quantity    uint64 `json:"quantity"`
	TotalPrice  uint64 `json:"total_price"`
}

type CartResponse struct {
	ID         uint64     `json:"id"`
	UserID     uint64     `json:"user_id"`
	GrantTotal uint64     `json:"grant_total"`
	CartItems  []CartItem `json:"cart_items"`
}

func FormatCart(cart *entities.CartModels) *CartResponse {
	cartResponse := &CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		GrantTotal: cart.GrandTotal,
	}

	var cartItems []CartItem
	for _, item := range cart.CartItems {
		cartItem := CartItem{
			CartItemID:  item.ID,
			ProductName: item.Product.Name,
			GramPlastic: item.Product.GramPlastic,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			TotalPrice:  item.TotalPrice,
		}
		cartItems = append(cartItems, cartItem)
	}

	cartResponse.CartItems = cartItems
	return cartResponse
}
