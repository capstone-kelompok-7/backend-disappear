package dto

type CreateCarouselRequest struct {
	Name  string `form:"name" validate:"required"`
	Photo string `form:"photo"`
}

type UpdateCarouselRequest struct {
	Name  string `form:"name"`
	Photo string `form:"photo"`
}
