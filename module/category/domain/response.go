package domain

type CategoryFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func FormatCategory(category *CategoryModels) *CategoryFormatter {
	categoryFormatter := &CategoryFormatter{}
	categoryFormatter.ID = category.ID
	categoryFormatter.Name = category.Name

	return categoryFormatter
}

func FormatterCategory(categories []*CategoryModels) []*CategoryFormatter {
	categoryFormatters := make([]*CategoryFormatter, 0)

	for _, category := range categories {
		formattedCategory := FormatCategory(category)
		categoryFormatters = append(categoryFormatters, formattedCategory)
	}

	return categoryFormatters
}
