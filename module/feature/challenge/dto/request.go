package dto

type CreateRequest struct {
	Title       string `form:"title" json:"title" validate:"required"`
	Photo       string `form:"photo" json:"photo" `
	StartDate   string `form:"start_date" json:"start_date" validate:"required"`
	EndDate     string `form:"end_date" json:"end_date" validate:"required"`
	Description string `form:"description" json:"description" validate:"required"`
	Exp         int    `form:"exp" json:"exp" validate:"required"`
}

type UpdateRequest struct {
	Title       string `form:"title" json:"title"`
	Photo       string `form:"photo" json:"photo"`
	StartDate   string `form:"start_date" json:"start_date"`
	EndDate     string `form:"end_date" json:"end_date"`
	Description string `form:"description" json:"description"`
	Exp         int    `form:"exp" json:"exp"`
}
