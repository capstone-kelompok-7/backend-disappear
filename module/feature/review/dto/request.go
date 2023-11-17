package dto

type CreateReviewRequest struct {
	UserID      uint64 `json:"user_id"`
	ProductID   uint64 `json:"product_id"`
	Rating      uint64 `json:"rating"`
	Description string `json:"description"`
}
