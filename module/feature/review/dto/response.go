package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

type ReviewFormatter struct {
	ID          uint64                  `json:"id"`
	UserID      uint64                  `json:"user_id"`
	ProductID   uint64                  `json:"product_id"`
	Rating      uint64                  `json:"rating"`
	Description string                  `json:"description"`
	Photos      []ReviewPhotosFormatter `json:"photos"`
}

type ReviewPhotosFormatter struct {
	ID       uint64 `json:"id"`
	ImageURL string `json:"url"`
}

func FormatReview(review *entities.ReviewModels) *ReviewFormatter {
	reviewFormatter := &ReviewFormatter{}
	reviewFormatter.ID = review.ID
	reviewFormatter.UserID = review.UserID
	reviewFormatter.ProductID = review.ProductID
	reviewFormatter.Rating = review.Rating
	reviewFormatter.Description = review.Description

	var image []ReviewPhotosFormatter
	for _, images := range review.Photos {
		categoryFormatter := ReviewPhotosFormatter{
			ID:       images.ID,
			ImageURL: images.ImageURL,
		}
		image = append(image, categoryFormatter)
	}
	reviewFormatter.Photos = image
	return reviewFormatter
}

func FormatReviewPhoto(reviewPhoto *entities.ReviewPhotoModels) *ReviewPhotosFormatter {
	photoFormatter := &ReviewPhotosFormatter{}
	photoFormatter.ID = reviewPhoto.ID
	photoFormatter.ImageURL = reviewPhoto.ImageURL
	return photoFormatter
}

func FormatterReview(reviews []*entities.ReviewModels) []*ReviewFormatter {
	reviewFormatters := make([]*ReviewFormatter, 0)

	for _, review := range reviews {
		formattedAddress := FormatReview(review)
		reviewFormatters = append(reviewFormatters, formattedAddress)
	}

	return reviewFormatters
}

type ReviewDetail struct {
	Name         string    `json:"name"`
	PhotoProfile string    `json:"photo_profile"`
	Rating       uint64    `json:"rating"`
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
}

func FormatReviewDetails(reviews []*entities.ReviewDetail) []*ReviewDetail {
	formattedReviews := make([]*ReviewDetail, 0)

	for _, review := range reviews {
		formattedReview := &ReviewDetail{
			Name:         review.Name,
			PhotoProfile: review.PhotoProfile,
			Rating:       review.Rating,
			Date:         review.Date,
			Description:  review.Description,
		}
		formattedReviews = append(formattedReviews, formattedReview)
	}

	return formattedReviews
}

// CreateReviewResponse formatter create review
type CreateReviewResponse struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	ProductID   uint64 `json:"product_id"`
	Rating      uint64 `json:"rating"`
	Description string `json:"description"`
}

func CreateFormatReview(review *entities.ReviewModels) *CreateReviewResponse {
	reviewFormatter := &CreateReviewResponse{}
	reviewFormatter.ID = review.ID
	reviewFormatter.UserID = review.UserID
	reviewFormatter.ProductID = review.ProductID
	reviewFormatter.Rating = review.Rating
	reviewFormatter.Description = review.Description

	return reviewFormatter
}

type ProductDetailFormatter struct {
	ID                 uint64                  `json:"id"`
	Name               string                  `json:"name"`
	CurrentRatingFive  uint64                  `json:"current_rating_five"`
	CurrentRatingFour  uint64                  `json:"current_rating_four"`
	CurrentRatingThree uint64                  `json:"current_rating_three"`
	CurrentRatingTwo   uint64                  `json:"current_rating_two"`
	CurrentRatingOne   uint64                  `json:"current_rating_one"`
	Rating             float64                 `json:"rating"`
	TotalReview        uint64                  `json:"total_review"`
	Reviews            []DetailReviewFormatter `json:"reviews"`
}

type DetailReviewFormatter struct {
	ID           uint64                `json:"id"`
	UserID       uint64                `json:"user_id"`
	Name         string                `json:"name"`
	PhotoProfile string                `json:"photo_profile"`
	Rating       uint64                `json:"rating"`
	Description  string                `json:"description"`
	Date         time.Time             `json:"date"`
	Photo        []ReviewPhotoResponse `json:"photo"`
}

type ReviewPhotoResponse struct {
	ID    uint64 `json:"id"`
	Photo string `json:"photo"`
}

func ConvertReviewPhotoModelsToResponse(photos []entities.ReviewPhotoModels) []ReviewPhotoResponse {
	var photoResponses []ReviewPhotoResponse

	for _, photo := range photos {
		photoResponse := ReviewPhotoResponse{
			ID:    photo.ID,
			Photo: photo.ImageURL,
		}

		photoResponses = append(photoResponses, photoResponse)
	}

	return photoResponses
}

func FormatProductDetail(product *entities.ProductModels) *ProductDetailFormatter {
	productFormatter := &ProductDetailFormatter{
		ID:          product.ID,
		Name:        product.Name,
		Rating:      product.Rating,
		TotalReview: product.TotalReview,
	}

	var currentRatingFive, currentRatingFour, currentRatingThree, currentRatingTwo, currentRatingOne uint64
	for _, review := range product.ProductReview {
		switch review.Rating {
		case 5:
			currentRatingFive++
		case 4:
			currentRatingFour++
		case 3:
			currentRatingThree++
		case 2:
			currentRatingTwo++
		case 1:
			currentRatingOne++
		}
	}

	productFormatter.CurrentRatingFive = currentRatingFive
	productFormatter.CurrentRatingFour = currentRatingFour
	productFormatter.CurrentRatingThree = currentRatingThree
	productFormatter.CurrentRatingTwo = currentRatingTwo
	productFormatter.CurrentRatingOne = currentRatingOne

	var reviews []DetailReviewFormatter
	for _, review := range product.ProductReview {
		reviewFormatter := DetailReviewFormatter{
			ID:           review.ID,
			UserID:       review.UserID,
			Name:         review.User.Name,
			PhotoProfile: review.User.PhotoProfile,
			Rating:       review.Rating,
			Description:  review.Description,
			Date:         review.Date,
			Photo:        ConvertReviewPhotoModelsToResponse(review.Photos),
		}
		reviews = append(reviews, reviewFormatter)

	}
	productFormatter.Reviews = reviews
	return productFormatter
}
