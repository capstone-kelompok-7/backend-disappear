package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type LoginResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type VerifyOTPResponse struct {
	AccessToken string `json:"access_token"`
}

// UserDetailResponse for detail users
type UserDetailResponse struct {
	ID             uint64 `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	PhotoProfile   string `json:"photo_profile"`
	TotalGram      uint64 `json:"total_gram"`
	TotalChallenge uint64 `json:"total_challenge"`
	IsVerified     bool   `json:"is_verified"`
	Level          string `json:"level"`
	Exp            uint64 `json:"exp"`
}

func FormatterDetailUser(user *entities.UserModels) *UserDetailResponse {
	userFormatter := &UserDetailResponse{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		Name:           user.Name,
		Phone:          user.Phone,
		PhotoProfile:   user.PhotoProfile,
		TotalGram:      user.TotalGram,
		TotalChallenge: user.TotalChallenge,
		IsVerified:     user.IsVerified,
		Level:          user.Level,
		Exp:            user.Exp,
	}
	return userFormatter
}
