package response

type PaymentMethodResponse struct {
	ID                 uint   `json:"id"`
	PaymentMethod      string `json:"paymentMethod"`
	BankCode           string `json:"bankCode"`
	BankName           string `json:"bankName"`
	AdminFee           int    `json:"adminFee"`
	MethodLogo         string `json:"methodLogo"`
	IsPercentage       bool   `json:"isPercentage"`
	AdminFeePercentage string `json:"adminFeePercentage"`
}
