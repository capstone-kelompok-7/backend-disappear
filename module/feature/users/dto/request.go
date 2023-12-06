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

type UserPersonalizationRequest struct {
	IsuID      []uint64 `json:"isu_id"`
	CategoryID []uint64 `json:"category_id"`
}

type UserPreferenceRequest struct {
	PreferredTopics []string `json:"preferred_topics" validate:"required"`
}
