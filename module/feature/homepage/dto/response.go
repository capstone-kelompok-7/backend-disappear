package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

type ContentResponse struct {
	Carousel []CarouselResponse `json:"carousel"`
	Category []CategoryResponse `json:"category"`
	Product  []ProductResponse  `json:"product"`
}

type CarouselResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type CategoryResponse struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	TotalProduct uint64 `json:"total_product"`
}

type ProductResponse struct {
	ID     uint64                `json:"id"`
	Name   string                `json:"name"`
	Rating float64               `json:"rating"`
	Price  uint64                `json:"price"`
	Photos *ProductPhotoResponse `json:"photos"`
}

type ProductPhotoResponse struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	URL       string `json:"url"`
}

func FormatContentResponse(carousels []*entities.CarouselModels, categories []*entities.CategoryModels, products []*entities.ProductModels) *ContentResponse {
	carouselResponses := make([]CarouselResponse, 0, len(carousels))
	for _, carousel := range carousels {
		carouselResponses = append(carouselResponses, CarouselResponse{
			ID:    carousel.ID,
			Name:  carousel.Name,
			Photo: carousel.Photo,
		})
	}

	categoryResponses := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		categoryResponses = append(categoryResponses, CategoryResponse{
			ID:           category.ID,
			Name:         category.Name,
			Photo:        category.Photo,
			TotalProduct: category.TotalProduct,
		})
	}

	productResponses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		productPhotos := &ProductPhotoResponse{}
		if len(product.ProductPhotos) > 0 {
			productPhotos.ID = product.ProductPhotos[0].ID
			productPhotos.ProductID = product.ProductPhotos[0].ProductID
			productPhotos.URL = product.ProductPhotos[0].ImageURL
		}

		productResponses = append(productResponses, ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Rating: product.Rating,
			Price:  product.Price,
			Photos: productPhotos,
		})
	}

	return &ContentResponse{
		Carousel: carouselResponses,
		Category: categoryResponses,
		Product:  productResponses,
	}
}

type BlogPostResponse struct {
	Challenge []ChallengeResponse `json:"challenge"`
	Articles  []ArticleResponse   `json:"articles"`
}

type ArticleResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Photo     string    `json:"photo"`
	Content   string    `json:"content"`
	Authors   string    `json:"authors"`
	Views     uint64    `json:"views"`
	CreatedAt time.Time `json:"created_at"`
}

type ChallengeResponse struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Photo string `json:"photo"`
	Exp   uint64 `json:"exp"`
}

func FormatBlogPostResponse(challenges []*entities.ChallengeModels, articles []*entities.ArticleModels) *BlogPostResponse {
	challengeResponses := make([]ChallengeResponse, 0, len(challenges))
	for _, challenge := range challenges {
		challengeResponses = append(challengeResponses, ChallengeResponse{
			ID:    challenge.ID,
			Title: challenge.Title,
			Photo: challenge.Photo,
			Exp:   challenge.Exp,
		})
	}

	articleResponses := make([]ArticleResponse, 0, len(articles))
	for _, article := range articles {
		articleResponses = append(articleResponses, ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Photo:     article.Photo,
			Content:   article.Content,
			Authors:   article.Author,
			Views:     article.Views,
			CreatedAt: article.CreatedAt,
		})
	}

	return &BlogPostResponse{
		Challenge: challengeResponses,
		Articles:  articleResponses,
	}
}
