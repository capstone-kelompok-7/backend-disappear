package dto

type CreateProductRequest struct {
	Name        string   `json:"name" form:"name" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	GramPlastic int      `json:"gram_plastic" form:"gram_plastic" validate:"required"`
	Price       float64  `json:"price" form:"price" validate:"required"`
	Stock       int      `json:"stock" form:"stock" validate:"required"`
	Discount    int      `json:"discount" form:"discount" validate:"required"`
	Exp         int      `json:"exp" form:"exp" validate:"required"`
	Categories  []uint64 `json:"categories" form:"categories" validate:"required"`
}
