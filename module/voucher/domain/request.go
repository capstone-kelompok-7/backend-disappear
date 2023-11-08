package domain

type VoucherRequestModel struct {
	ID          uint64  `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Category    string  `json:"category" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Discouunt   int     `json:"discount" validate:"required"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     string  `json:"end_date" validate:"required"`
	MinAmount   float64 `json:"min_amount" validate:"required"`
}
