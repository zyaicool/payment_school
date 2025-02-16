package models

import "time"

type BillingDetail struct {
	Master
	BillingID         uint       `json:"billingId"`
	DueDate           *time.Time `json:"dueDate"`
	DetailBillingName string     `json:"detailBillingName"`
	Amount            int64      `json:"amount"`
}
