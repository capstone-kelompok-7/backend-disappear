package sendnotif

type SendNotificationRequest struct {
	OrderID string `json:"order_id"`
	UserID  uint64 `json:"user_id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Token   string `json:"token"`
}
