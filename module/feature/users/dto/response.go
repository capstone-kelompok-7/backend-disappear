package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

// UserDetailResponse for detail users
type UserDetailResponse struct {
	ID             uint64            `json:"id"`
	Email          string            `json:"email"`
	Role           string            `json:"role"`
	Name           string            `json:"name"`
	Phone          string            `json:"phone"`
	PhotoProfile   string            `json:"photo_profile"`
	TotalGram      uint64            `json:"total_gram"`
	TotalChallenge uint64            `json:"total_challenge"`
	IsVerified     bool              `json:"is_verified"`
	Level          string            `json:"level"`
	Exp            uint64            `json:"exp"`
	Addresses      []AddressResponse `json:"addresses"`
}

type AddressResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	AcceptedName string `json:"accepted_name" `
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	IsPrimary    bool   `json:"is_primary"`
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
	var addresses []AddressResponse
	for _, address := range user.Address {
		addressesFormatter := AddressResponse{
			ID:           address.ID,
			AcceptedName: address.AcceptedName,
			Phone:        address.Phone,
			Address:      address.Address,
			IsPrimary:    address.IsPrimary,
		}
		addresses = append(addresses, addressesFormatter)
	}
	userFormatter.Addresses = addresses

	return userFormatter
}

// UserPaginationResponse for get all pagination
type UserPaginationResponse struct {
	ID             uint64 `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	PhotoProfile   string `json:"photo_profile"`
	TotalGram      uint64 `json:"total_gram"`
	TotalChallenge uint64 `json:"total_challenge"`
	Level          string `json:"level"`
	Exp            uint64 `json:"exp"`
}

func FormatUserPagination(user *entities.UserModels) *UserPaginationResponse {
	userFormatter := &UserPaginationResponse{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		Name:           user.Name,
		Phone:          user.Phone,
		PhotoProfile:   user.PhotoProfile,
		TotalGram:      user.TotalGram,
		TotalChallenge: user.TotalChallenge,
		Level:          user.Level,
		Exp:            user.Exp,
	}
	return userFormatter
}

func FormatterUsersPagination(users []*entities.UserModels) []*UserPaginationResponse {
	usersFormatters := make([]*UserPaginationResponse, 0)

	for _, user := range users {
		formattedUsers := FormatUserPagination(user)
		usersFormatters = append(usersFormatters, formattedUsers)
	}

	return usersFormatters
}

// UserLeaderboardResponse for get leaderboard
type UserLeaderboardResponse struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name"`
	PhotoProfile   string `json:"photo_profile"`
	TotalGram      uint64 `json:"total_gram"`
	TotalChallenge uint64 `json:"total_challenge"`
	Level          string `json:"level"`
	Exp            uint64 `json:"exp"`
}

func FormatUserLeaderboard(user *entities.UserModels) *UserLeaderboardResponse {
	userFormatter := &UserLeaderboardResponse{
		ID:             user.ID,
		Name:           user.Name,
		PhotoProfile:   user.PhotoProfile,
		TotalGram:      user.TotalGram,
		TotalChallenge: user.TotalChallenge,
		Level:          user.Level,
		Exp:            user.Exp,
	}
	return userFormatter
}

func FormatterUserLeaderboard(users []*entities.UserModels) []*UserLeaderboardResponse {
	usersFormatters := make([]*UserLeaderboardResponse, 0)

	for _, user := range users {
		formattedUsers := FormatUserLeaderboard(user)
		usersFormatters = append(usersFormatters, formattedUsers)
	}

	return usersFormatters
}

// UserActivityResponse for user activities
type UserActivityResponse struct {
	NumSuccessfulOrders int `json:"num_successful_orders"`
	NumFailedOrders     int `json:"num_failed_orders"`
	TotalOrders         int `json:"total_orders"`
	NumSuccessfulChall  int `json:"num_successful_challenges"`
	NumFailedChall      int `json:"num_failed_challenges"`
	TotalChallenges     int `json:"total_challenges"`
}

func FormatUserActivityResponse(numSuccessfulOrders, numFailedOrders, totalOrders, numSuccessfulChallenge, numFailedChallenge, totalChallenge int) *UserActivityResponse {
	return &UserActivityResponse{
		NumSuccessfulOrders: numSuccessfulOrders,
		NumFailedOrders:     numFailedOrders,
		TotalOrders:         totalOrders,
		NumSuccessfulChall:  numSuccessfulChallenge,
		NumFailedChall:      numFailedChallenge,
		TotalChallenges:     totalChallenge,
	}
}

// UserProfileResponse for get profile
type UserProfileResponse struct {
	ID             uint64 `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	PhotoProfile   string `json:"photo_profile"`
	TotalGram      uint64 `json:"total_gram"`
	TotalChallenge uint64 `json:"total_challenge"`
	Level          string `json:"level"`
}

func FormatUserProfileResponse(user *entities.UserModels) *UserProfileResponse {
	userFormatter := &UserProfileResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		Phone:          user.Phone,
		PhotoProfile:   user.PhotoProfile,
		TotalGram:      user.TotalGram,
		TotalChallenge: user.TotalChallenge,
		Level:          user.Level,
	}
	return userFormatter
}
