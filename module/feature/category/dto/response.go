package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type CategoryFormatter struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	TotalProduct uint64 `json:"total_product"`
}

func FormatCategory(category *entities.CategoryModels) *CategoryFormatter {
	categoryFormatter := &CategoryFormatter{}
	categoryFormatter.ID = category.ID
	categoryFormatter.Name = category.Name
	categoryFormatter.Photo = category.Photo
	categoryFormatter.TotalProduct = category.TotalProduct

	return categoryFormatter
}

func FormatterCategory(categories []*entities.CategoryModels) []*CategoryFormatter {
	categoryFormatters := make([]*CategoryFormatter, 0)

	for _, category := range categories {
		formattedCategory := FormatCategory(category)
		categoryFormatters = append(categoryFormatters, formattedCategory)
	}

	return categoryFormatters
}
