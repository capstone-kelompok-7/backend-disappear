package dto

type CardResponse struct {
	ProductCount int64   `json:"product_count"`
	UserCount    int64   `json:"user_count"`
	OrderCount   int64   `json:"order_count"`
	IncomeCount  float64 `json:"income_count"`
}

func FormatCardResponse(productCount, userCount, orderCount int64, inComeCount float64) *CardResponse {
	return &CardResponse{
		ProductCount: productCount,
		UserCount:    userCount,
		OrderCount:   orderCount,
		IncomeCount:  inComeCount,
	}
}
