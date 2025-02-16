package models

type PaymentMethod struct {
	Master
	ID                 uint   `json:"id" gorm:"primaryKey"`
	PaymentMethod      string `json:"paymentMethod"`
	BankCode           string `json:"bankCode"`
	BankName           string `json:"bankName"`
	AdminFee           int    `json:"adminFee"`
	MethodLogo         string `json:"methodLogo"`
	IsPercentage       bool   `json:"isPercentage"`
	AdminFeePercentage string `json:"adminFeePercentage"`
}

// TableName overrides the table name used by GORM
func (PaymentMethod) TableName() string {
	return "master_payment_method"
}
