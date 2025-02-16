package models

import "time"

type TransactionBillingDetail struct {
	Master
	TransactionBillingID  uint       `json:"transactionBillingId"`
	MasterPaymentMethodID uint       `json:"masterPaymentMethodId"`
	Discount              int        `json:"discount"`
	DiscountType          string     `json:"discountType"`
	ChangeAmount          int        `json:"changeAmount"`
	BankName              *string    `gorm:"column:bank_name"`
	VirtualAccountNumber  *string    `gorm:"column:virtual_account_number"`
	TransactionTime       *time.Time `json:"transactionTime"`
	IsDonation            bool       `json:"isDonation" gorm:"default:false"`
}
