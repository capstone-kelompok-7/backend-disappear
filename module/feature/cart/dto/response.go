package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type CartItemFormatter struct {
	CartItemID  uint64          `json:"cart_item_id"`
	ProductName string          `json:"product_name"`
	GramPlastic uint64          `json:"gram_plastic"`
	Price       uint64          `json:"price"`
	Quantity    uint64          `json:"quantity"`
	TotalPrice  uint64          `json:"total_price"`
	Product     ProductResponse `json:"product,omitempty"`
}

type CartFormatter struct {
	ID         uint64              `json:"id"`
	UserID     uint64              `json:"user_id"`
	GrantTotal uint64              `json:"grant_total"`
	CartItems  []CartItemFormatter `json:"cart_items"`
}

type ProductPhotoResponse struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	URL       string `json:"url"`
}

type ProductResponse struct {
	ID            uint64                 `json:"id"`
	Name          string                 `json:"name"`
	Price         uint64                 `json:"price"`
	Discount      uint64                 `json:"discount"`
	ProductPhotos []ProductPhotoResponse `json:"product_photos"`
}

func FormatCart(cart *entities.CartModels) *CartFormatter {
	cartResponse := &CartFormatter{
		ID:         cart.ID,
		UserID:     cart.UserID,
		GrantTotal: cart.GrandTotal,
	}

	var cartItems []CartItemFormatter
	for _, item := range cart.CartItems {
		var productPhotos []ProductPhotoResponse
		for _, photo := range item.Product.ProductPhotos {
			productPhotos = append(productPhotos, ProductPhotoResponse{
				ID:        photo.ID,
				ProductID: photo.ProductID,
				URL:       photo.ImageURL,
			})
		}

		cartItem := CartItemFormatter{
			CartItemID:  item.ID,
			ProductName: item.Product.Name,
			GramPlastic: item.Product.GramPlastic,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			TotalPrice:  item.TotalPrice,
			Product: ProductResponse{
				ID:            item.Product.ID,
				Name:          item.Product.Name,
				Price:         item.Product.Price,
				Discount:      item.Product.Discount,
				ProductPhotos: productPhotos,
			},
		}
		if len(item.Product.ProductPhotos) > 0 {
			productPhoto := ProductPhotoResponse{
				ID:        item.Product.ProductPhotos[0].ID,
				ProductID: item.Product.ProductPhotos[0].ProductID,
				URL:       item.Product.ProductPhotos[0].ImageURL,
			}
			cartItem.Product.ProductPhotos = []ProductPhotoResponse{productPhoto}
		}
		cartItems = append(cartItems, cartItem)
	}

	cartResponse.CartItems = cartItems
	return cartResponse
}
