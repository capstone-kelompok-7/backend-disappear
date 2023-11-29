package dto

type CreateReviewRequest struct {
	UserID      uint64 `json:"user_id"`
	ProductID   uint64 `json:"product_id" validate:"required"`
	Rating      uint64 `json:"rating" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreatePhotoReviewRequest struct {
	ReviewID uint64 `form:"review_id" json:"review_id"`
	Photo    string `form:"photo" json:"photo"`
}
