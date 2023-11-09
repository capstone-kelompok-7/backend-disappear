package dto

type RegisterRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type EmailRequest struct {
	Email string `form:"email" json:"email" validate:"required"`
	OTP   string `form:"otp" json:"otp" validate:"required"`
}

type ResendOTPRequest struct {
	Email string `form:"email" json:"email" validate:"required"`
}
