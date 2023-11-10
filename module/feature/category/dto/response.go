package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type CategoryFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func FormatCategory(category *entities.CategoryModels) *CategoryFormatter {
	categoryFormatter := &CategoryFormatter{}
	categoryFormatter.ID = category.ID
	categoryFormatter.Name = category.Name

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
