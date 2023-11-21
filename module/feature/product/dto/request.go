package dto

type CreateProductRequest struct {
	Name        string   `json:"name" form:"name" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	GramPlastic uint64   `json:"gram_plastic" form:"gram_plastic" validate:"required"`
	Price       uint64   `json:"price" form:"price" validate:"required"`
	Stock       uint64   `json:"stock" form:"stock" validate:"required"`
	Discount    uint64   `json:"discount" form:"discount" validate:"required"`
	Exp         uint64   `json:"exp" form:"exp" validate:"required"`
	Categories  []uint64 `json:"categories" form:"categories" validate:"required"`
	ImageURL    string   `json:"image_url" form:"image_url"`
}

type CreateProductImage struct {
	ProductID uint64 `form:"product_id"`
	Image     string `form:"photo" validate:"required"`
}

type UpdateProduct struct {
	Name        string   `json:"name" form:"name" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	GramPlastic uint64   `json:"gram_plastic" form:"gram_plastic" validate:"required"`
	Price       uint64   `json:"price" form:"price" validate:"required"`
	Stock       uint64   `json:"stock" form:"stock" validate:"required"`
	Discount    uint64   `json:"discount" form:"discount" validate:"required"`
	Exp         uint64   `json:"exp" form:"exp" validate:"required"`
	CategoryIDs []uint64 `json:"categories" form:"categories" validate:"required"`
	ImageURL    string   `json:"image_url" form:"image_url"`
}
