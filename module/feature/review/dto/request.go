package dto

type CreateReviewRequest struct {
	UserID      uint64 `json:"user_id"`
	ProductID   uint64 `json:"product_id"`
	Rating      uint64 `json:"rating"`
	Description string `json:"description"`
}

type CreatePhotoReviewRequest struct {
	ReviewID uint64 `form:"review_id" json:"review_id"`
	Photo    string `form:"photo" json:"photo"`
}
