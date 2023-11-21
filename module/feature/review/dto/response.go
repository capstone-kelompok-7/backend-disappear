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
