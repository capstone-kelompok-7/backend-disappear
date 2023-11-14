package dto

type CreateCategoryRequest struct {
	Name  string `form:"name" validate:"required"`
	Photo string `form:"photo"`
}

type UpdateCategoryRequest struct {
	Name  string `form:"name"`
	Photo string `form:"photo"`
}
