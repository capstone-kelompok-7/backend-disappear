package dto

type RegisterRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6,noSpace"`
}

type LoginRequest struct {
	Email       string `form:"email" json:"email" validate:"required,email"`
	Password    string `form:"password" json:"password" validate:"required,min=6,noSpace"`
	DeviceToken string `form:"device_token" json:"device_token"`
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

type RegisterSocialRequest struct {
	SocialID     string `json:"social_id" validate:"required"`
	Provider     string `json:"provider"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	PhotoProfile string `json:"photo_profile"`
}

type LoginSocialRequest struct {
	SocialID string `json:"social_id" validate:"required"`
}
