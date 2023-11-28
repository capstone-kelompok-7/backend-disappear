package dto

type CreateCarouselRequest struct {
	Name  string `form:"name" json:"name" validate:"required"`
	Photo string `form:"photo" json:"photo"`
}

type UpdateCarouselRequest struct {
	Name  string `form:"name" json:"name"`
	Photo string `form:"photo" json:"photo"`
}
