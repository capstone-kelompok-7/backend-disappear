package homepage

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/labstack/echo/v4"
)

type RepositoryHomepageInterface interface {
	GetBestSellingProducts(limit int) ([]*entities.ProductModels, error)
	GetFiveCategories() ([]*entities.CategoryModels, error)
	GetFiveCarousel() ([]*entities.CarouselModels, error)
	GetFiveChallenge() ([]*entities.ChallengeModels, error)
	GetThreeArticle() ([]*entities.ArticleModels, error)
}

type ServiceHomepageInterface interface {
	GetBestSellingProducts(limit int) ([]*entities.ProductModels, error)
	GetCategory() ([]*entities.CategoryModels, error)
	GetCarousel() ([]*entities.CarouselModels, error)
	GetChallenge() ([]*entities.ChallengeModels, error)
	GetArticle() ([]*entities.ArticleModels, error)
}

type HandlerHomepageInterface interface {
	GetHomepageContent() echo.HandlerFunc
	GetBlogPost() echo.HandlerFunc
}
