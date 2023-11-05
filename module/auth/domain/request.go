package domain

type RegisterRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Phone    string `form:"phone" json:"phone" validate:"required"`
	Password string `form:"password" json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}
