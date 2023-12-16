package dto

type FcmRequest struct {
	Title  string `json:"title" form:"title"`
	UserID uint64 `json:"user_id" form:"user_id"`
	Body   string `json:"body" form:"body"`
	Token  string `json:"token" form:"token"`
}
