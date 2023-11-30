package dto

type CreateAddressRequest struct {
	UserID       uint64 `form:"user_id" json:"user_id"`
	AcceptedName string `form:"accepted_name" json:"accepted_name" validate:"required"`
	Phone        string `form:"phone" json:"phone" validate:"required"`
	Address      string `form:"address" json:"address" validate:"required"`
	IsPrimary    bool   `form:"is_primary" json:"is_primary"`
}

type UpdateAddressRequest struct {
	UserID       uint64 `form:"user_id" json:"user_id"`
	AcceptedName string `form:"accepted_name" json:"accepted_name"`
	Phone        string `form:"phone" json:"phone"`
	Address      string `form:"address" json:"address"`
	IsPrimary    bool   `form:"is_primary" json:"is_primary"`
}
