package dto

type LoginResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}
