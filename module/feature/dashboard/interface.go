package dashboard

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
	"time"
)

type RepositoryDashboardInterface interface {
	CountProducts() (int64, error)
	CountUsers() (int64, error)
	CountOrder() (int64, error)
	CountIncome() (float64, error)
	CountTotalGram() (int64, error)
	GetProductWithMaxReviews() ([]*entities.ProductModels, error)
	GetGramPlasticStat(startOfWeek, endOfWeek time.Time) (uint64, error)
}

type ServiceDashboardInterface interface {
	GetCardDashboard() (int64, int64, int64, float64, error)
	GetLandingPage() (int64, int64, int64, error)
	GetProductReviewsWithMaxTotal() ([]*entities.ProductModels, error)
	GetGramPlasticStat(startOfWeek, endOfWeek time.Time) (uint64, error)
}

type HandlerDashboardInterface interface {
	GetCardDashboard() echo.HandlerFunc
	GetLandingPage() echo.HandlerFunc
	GetReview() echo.HandlerFunc
	GetGramPlasticStat() echo.HandlerFunc
}
