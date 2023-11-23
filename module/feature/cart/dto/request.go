package dto

type AddCartItemsRequest struct {
	UserID    uint64 `form:"user_id" json:"user_id"`
	ProductID uint64 `form:"product_id" json:"product_id" validate:"required"`
	Quantity  uint64 `form:"quantity" json:"quantity" validate:"required"`
}

type ReduceCartItemsRequest struct {
	CartItemID uint64 `form:"cart_item_id" json:"cart_item_id"`
	Quantity   uint64 `form:"quantity" json:"quantity" validate:"required"`
}
