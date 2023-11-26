package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type ProductFormatter struct {
	ID          uint64                  `json:"id"`
	Name        string                  `json:"name"`
	GramPlastic uint64                  `json:"gram_plastic"`
	Stock       uint64                  `json:"stock"`
	Discount    uint64                  `json:"discount"`
	Exp         uint64                  `json:"product_exp"`
	Price       uint64                  `json:"price"`
	Rating      float64                 `json:"rating"`
	TotalReview uint64                  `json:"total_review"`
	Categories  []CategoryFormatter     `json:"categories"`
	Images      []ProductImageFormatter `json:"image_url"`
	Reviews     []ReviewFormatter       `json:"reviews"`
}

type CategoryFormatter struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
type ProductImageFormatter struct {
	ID  uint64 `json:"id"`
	URL string `json:"image_url"`
}

type ReviewFormatter struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	Rating      uint64 `json:"rating"`
	Description string `json:"description"`
}
type CreateImageFormatter struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	Image     string `json:"image_url"`
}

type FormatProductCreated struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	GramPlastic uint64 `json:"gram_plastic"`
	Stock       uint64 `json:"stock"`
	Discount    uint64 `json:"discount"`
	Exp         uint64 `json:"product_exp"`
	Price       uint64 `json:"price"`
}

func ProductPhotoCreatedResponse(productPhoto *entities.ProductPhotosModels) CreateImageFormatter {
	response := CreateImageFormatter{}
	response.ID = productPhoto.ID
	response.ProductID = productPhoto.ProductID
	response.Image = productPhoto.ImageURL
	return response
}

func ProductCreatedResponse(product *entities.ProductModels) FormatProductCreated {
	productResponse := FormatProductCreated{}
	productResponse.ID = product.ID
	productResponse.Name = product.Name
	productResponse.GramPlastic = product.GramPlastic
	productResponse.Stock = product.Stock
	productResponse.Discount = product.Discount
	productResponse.Exp = product.Exp
	productResponse.Price = product.Price

	return productResponse
}

func FormatProduct(product *entities.ProductModels) *ProductFormatter {
	productFormatter := &ProductFormatter{}
	productFormatter.ID = product.ID
	productFormatter.Name = product.Name
	productFormatter.GramPlastic = product.GramPlastic
	productFormatter.Price = product.Price
	productFormatter.Stock = product.Stock
	productFormatter.Discount = product.Discount
	productFormatter.Exp = product.Exp
	productFormatter.Rating = product.Rating
	productFormatter.TotalReview = product.TotalReview

	var categories []CategoryFormatter
	for _, category := range product.Categories {
		categoryFormatter := CategoryFormatter{
			ID:   category.ID,
			Name: category.Name,
		}
		categories = append(categories, categoryFormatter)
	}
	productFormatter.Categories = categories

	var images []ProductImageFormatter
	for _, photo := range product.ProductPhotos {
		if photo.DeletedAt == nil {
			image := ProductImageFormatter{
				ID:  photo.ID,
				URL: photo.ImageURL,
			}
			images = append(images, image)
		}
	}
	productFormatter.Images = images

	var reviews []ReviewFormatter
	for _, review := range product.ProductReview {
		reviewFormatter := ReviewFormatter{
			ID:          review.ID,
			UserID:      review.UserID,
			Rating:      review.Rating,
			Description: review.Description,
		}
		reviews = append(reviews, reviewFormatter)
	}
	productFormatter.Reviews = reviews

	return productFormatter
}

func FormatterProduct(products []*entities.ProductModels) []*ProductFormatter {
	var productFormatter []*ProductFormatter

	for _, product := range products {
		formatProduct := FormatProduct(product)
		productFormatter = append(productFormatter, formatProduct)
	}

	return productFormatter
}

type ProductDetailFormatter struct {
	ID          uint64                  `json:"id"`
	Name        string                  `json:"name"`
	GramPlastic uint64                  `json:"gram_plastic"`
	Stock       uint64                  `json:"stock"`
	Discount    uint64                  `json:"discount"`
	Exp         uint64                  `json:"product_exp"`
	Price       uint64                  `json:"price"`
	Categories  []CategoryFormatter     `json:"categories,omitempty"`
	Images      []ProductImageFormatter `json:"image_url,omitempty"`
}

func FormatProductDetail(product entities.ProductModels) ProductDetailFormatter {
	productFormatter := ProductDetailFormatter{
		ID:          product.ID,
		Name:        product.Name,
		GramPlastic: product.GramPlastic,
		Price:       product.Price,
		Stock:       product.Stock,
		Discount:    product.Discount,
		Exp:         product.Exp,
	}

	var categories []CategoryFormatter
	for _, category := range product.Categories {
		categoryFormatter := CategoryFormatter{
			ID:   category.ID,
			Name: category.Name,
		}
		categories = append(categories, categoryFormatter)
	}
	productFormatter.Categories = categories

	var images []ProductImageFormatter
	for _, photo := range product.ProductPhotos {
		image := ProductImageFormatter{
			ID:  photo.ID,
			URL: photo.ImageURL,
		}
		images = append(images, image)
	}
	productFormatter.Images = images

	return productFormatter
}

type ReviewProductFormatter struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Rating      float64 `json:"rating"`
	TotalReview uint64  `json:"total_review"`
}

func FormatReviewProductFormatter(products []*entities.ProductModels) []*ReviewProductFormatter {
	productFormatters := make([]*ReviewProductFormatter, 0)
	for _, product := range products {
		productFormatter := &ReviewProductFormatter{
			ID:          product.ID,
			Name:        product.Name,
			Rating:      product.Rating,
			TotalReview: product.TotalReview,
		}
		productFormatters = append(productFormatters, productFormatter)
	}
	return productFormatters
}
