package dto

type CreateAddressRequest struct {
	UserID       uint64 `form:"user_id" json:"user_id"`
	AcceptedName string `form:"accepted_name" json:"accepted_name" validate:"required"`
	Street       string `form:"street" json:"street" validate:"required"`
	SubDistrict  string `form:"sub_district" json:"sub_district" validate:"required"`
	City         string `form:"city" json:"city" validate:"required"`
	Province     string `form:"province" json:"province" validate:"required"`
	PostalCode   int    `form:"postal_code" json:"postal_code" validate:"required"`
	Note         string `form:"note" json:"note"`
}
