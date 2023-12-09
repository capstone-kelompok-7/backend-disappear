package dto

type FcmRequest struct {
	Title  string `json:"title" form:"title" validate:required"`
	UserID uint64 `json:"user_id" form:"user_id" validate:required`
	Body   string `json:"body" form:"body" validate:required"`
	Token  string `json:"token" form:"token" validate:required"`
}
