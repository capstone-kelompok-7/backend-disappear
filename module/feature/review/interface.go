package review

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryReviewInterface interface {
	CreateReview(newData *entities.ReviewModels) (*entities.ReviewModels, error)
}

type ServiceReviewInterface interface {
	CreateReview(newData *entities.ReviewModels) (*entities.ReviewModels, error)
}

type HandlerReviewInterface interface {
	CreateReview() echo.HandlerFunc
}
