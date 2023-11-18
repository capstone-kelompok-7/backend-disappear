package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type ProductFormatter struct {
	ID          uint64                  `json:"id"`
	Name        string                  `json:"name"`
	GramPlastic uint64                  `json:"gram_plastic"`
	Stock       uint64                  `json:"stock"`
	Discount    uint64                  `json:"discount"`
	Exp         uint64                  `json:"product_exp"`
	Price       float64                 `json:"price"`
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

func FormatProduct(product entities.ProductModels) ProductFormatter {
	productFormatter := ProductFormatter{}
	productFormatter.ID = product.ID
	productFormatter.Name = product.Name
	productFormatter.GramPlastic = product.GramPlastic
	productFormatter.Price = product.Price
	productFormatter.Stock = product.Stock
	productFormatter.Discount = product.Discount
	productFormatter.Exp = product.Exp

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

func FormatterProduct(products []entities.ProductModels) []ProductFormatter {
	var productFormatter []ProductFormatter

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
	Price       float64                 `json:"price"`
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
