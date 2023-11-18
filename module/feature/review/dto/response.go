package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type ReviewFormatter struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	ProductID   uint64 `json:"product_id"`
	Rating      uint64 `json:"rating"`
	Description string `json:"description"`
}

func FormatReview(address *entities.ReviewModels) *ReviewFormatter {
	reviewFormatter := &ReviewFormatter{}
	reviewFormatter.ID = address.ID
	reviewFormatter.UserID = address.UserID
	reviewFormatter.ProductID = address.ProductID
	reviewFormatter.Rating = address.Rating
	reviewFormatter.Description = address.Description

	return reviewFormatter
}

func FormatterAddress(reviews []*entities.ReviewModels) []*ReviewFormatter {
	reviewFormatters := make([]*ReviewFormatter, 0)

	for _, review := range reviews {
		formattedAddress := FormatReview(review)
		reviewFormatters = append(reviewFormatters, formattedAddress)
	}

	return reviewFormatters
}
