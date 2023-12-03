package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

type CardResponse struct {
	ProductCount int64   `json:"product_count"`
	UserCount    int64   `json:"user_count"`
	OrderCount   int64   `json:"order_count"`
	IncomeCount  float64 `json:"income_count"`
}

func FormatCardResponse(productCount, userCount, orderCount int64, inComeCount float64) *CardResponse {
	return &CardResponse{
		ProductCount: productCount,
		UserCount:    userCount,
		OrderCount:   orderCount,
		IncomeCount:  inComeCount,
	}
}

type LandingPageResponse struct {
	UserCount   int64 `json:"user_count"`
	GramPlastic int64 `json:"gram_plastic"`
	OrderCount  int64 `json:"order_count"`
}

func FormatLandingPage(userCount, gramPlastic, orderCount int64) *LandingPageResponse {
	return &LandingPageResponse{
		UserCount:   userCount,
		GramPlastic: gramPlastic,
		OrderCount:  orderCount,
	}
}

type LandingPageReviewResponse struct {
	ID     uint64            `json:"id"`
	Name   string            `json:"name"`
	Review []*ReviewResponse `json:"review"`
}

type ReviewResponse struct {
	ID          uint64        `json:"id"`
	UserID      uint64        `json:"user_id"`
	Rating      uint64        `json:"rating"`
	Description string        `json:"description"`
	Date        time.Time     `json:"date"`
	User        *UserResponse `json:"user"`
}

type UserResponse struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	PhotoProfile string `json:"photo_profile"`
}

func FormatLandingPageReview(products []*entities.ProductModels) []*LandingPageReviewResponse {
	var response []*LandingPageReviewResponse

	for _, product := range products {
		reviews := make([]*ReviewResponse, len(product.ProductReview))
		for i, review := range product.ProductReview {
			reviews[i] = &ReviewResponse{
				ID:          review.ID,
				UserID:      review.UserID,
				Rating:      review.Rating,
				Description: review.Description,
				Date:        review.Date,
				User: &UserResponse{
					ID:           review.User.ID,
					Name:         review.User.Name,
					PhotoProfile: review.User.PhotoProfile,
				},
			}
		}

		response = append(response, &LandingPageReviewResponse{
			ID:     product.ID,
			Name:   product.Name,
			Review: reviews,
		})
	}

	return response
}

type GramPlasticStat struct {
	Week           string `json:"week"`
	GramTotalCount uint64 `json:"gram_total_count"`
}

var MonthMap = map[time.Month]string{
	time.January:   "Januari",
	time.February:  "Februari",
	time.March:     "Maret",
	time.April:     "April",
	time.May:       "Mei",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "Agustus",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Desember",
}

type LastTransactionResponse struct {
	Username      string `json:"username"`
	TotalPrice    uint64 `json:"total_price"`
	Date          string `json:"date"`
	PaymentStatus string `json:"payment_status"`
}
