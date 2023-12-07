package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
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
	Description string                  `json:"description"`
	GramPlastic uint64                  `json:"gram_plastic"`
	Stock       uint64                  `json:"stock"`
	Discount    uint64                  `json:"discount"`
	Exp         uint64                  `json:"exp"`
	Price       uint64                  `json:"price"`
	Rating      float64                 `json:"rating"`
	TotalReview uint64                  `json:"total_review"`
	TotalSold   uint64                  `json:"total_sold"`
	Categories  []CategoryFormatter     `json:"categories"`
	Images      []ProductImageFormatter `json:"image_url"`
	Reviews     []ReviewFormatter       `json:"reviews"`
}

type CreateProductFormatter struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GramPlastic uint64 `json:"gram_plastic"`
	Stock       uint64 `json:"stock"`
	Discount    uint64 `json:"discount"`
	Exp         uint64 `json:"exp"`
	Price       uint64 `json:"price"`
}

func FormatCreateProductResponse(product *entities.ProductModels) CreateProductFormatter {
	createdProduct := CreateProductFormatter{}
	createdProduct.ID = product.ID
	createdProduct.Name = product.Name
	createdProduct.Description = product.Description
	createdProduct.GramPlastic = product.GramPlastic
	createdProduct.Stock = product.Stock
	createdProduct.Discount = product.Discount
	createdProduct.Exp = product.Exp
	createdProduct.Price = product.Price

	return createdProduct

}

type CreateImageFormatter struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	Image     string `json:"image_url"`
}

func ProductPhotoCreatedResponse(productPhoto *entities.ProductPhotosModels) CreateImageFormatter {
	response := CreateImageFormatter{}
	response.ID = productPhoto.ID
	response.ProductID = productPhoto.ProductID
	response.Image = productPhoto.ImageURL
	return response
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

func FormatProductDetail(product entities.ProductModels, totalSold uint64) ProductDetailFormatter {
	productFormatter := ProductDetailFormatter{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		GramPlastic: product.GramPlastic,
		Price:       product.Price,
		Stock:       product.Stock,
		Discount:    product.Discount,
		Exp:         product.Exp,
		Rating:      product.Rating,
		TotalReview: product.TotalReview,
		TotalSold:   totalSold,
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

	var reviews []ReviewFormatter
	for _, review := range product.ProductReview {
		reviewFormatter := ReviewFormatter{
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

type OtherProductFormatter struct {
	ID     uint64                  `json:"id"`
	Name   string                  `json:"name"`
	Price  uint64                  `json:"price"`
	Rating float64                 `json:"rating"`
	Images []ProductImageFormatter `json:"image_url"`
}

func FormatOtherProduct(product *entities.ProductModels) *OtherProductFormatter {
	productFormatter := &OtherProductFormatter{
		ID:     product.ID,
		Name:   product.Name,
		Price:  product.Price,
		Rating: product.Rating,
	}
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

	return productFormatter
}

func FormatterOtherProduct(products []*entities.ProductModels) []*OtherProductFormatter {
	var productFormatter []*OtherProductFormatter

	for _, product := range products {
		formatProduct := FormatOtherProduct(product)
		productFormatter = append(productFormatter, formatProduct)
	}

	return productFormatter
}
