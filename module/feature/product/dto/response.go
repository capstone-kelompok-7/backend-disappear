package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type ProductFormatter struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	GramPlastic int                 `json:"gram_plastic"`
	Stock       int                 `json:"stock"`
	Discount    int                 `json:"discount"`
	Exp         int                 `json:"product_exp"`
	Price       float64             `json:"price"`
	Categories  []CategoryFormatter `json:"categories"`
	Images      interface{}         `json:"image_url"`
}

type CategoryFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

	var images []map[string]interface{}
	for _, photo := range product.ProductPhotos {
		image := map[string]interface{}{
			"id":  photo.ID,
			"url": photo.ImageURL,
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
