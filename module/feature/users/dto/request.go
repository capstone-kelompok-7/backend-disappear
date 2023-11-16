package dto

type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type EditProfileRequest struct {
	Name         string `form:"name" json:"name"`
	Phone        string `form:"phone" json:"phone"`
	PhotoProfile string `form:"photo" json:"photo_profile"`
}
