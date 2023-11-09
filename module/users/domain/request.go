package domain

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	ConPassword string `json:"con_password"`
}
