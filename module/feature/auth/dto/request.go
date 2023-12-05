package dto

type RegisterRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6,noSpace"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6,noSpace"`
}

type EmailRequest struct {
	Email string `form:"email" json:"email" validate:"required,email"`
	OTP   string `form:"otp" json:"otp" validate:"required"`
}

type ResendOTPRequest struct {
	Email string `form:"email" json:"email" validate:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `form:"email" json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=6,noSpace"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,noSpace"`
}
