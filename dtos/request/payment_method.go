package request

type PaymentMethodCreateRequest struct {
	PaymentMethod      string `validate:"required"`
	BankCode           string `validate:"required"`
	BankName           string `validate:"required"`
	AdminFee           int    `validate:"required"`
	MethodLogo         string `validate:"omitempty"`
	IsPercentage       string `json:"isPercentage"`
	AdminFeePercentage string `json:"adminFeePercentage"`
}
