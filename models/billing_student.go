package models

import "time"

type BillingStudent struct {
	Master
	BillingID         uint       `json:"billingId"`
	StudentID         uint       `json:"studentId"`
	PaymentStatus     string     `json:"paymentStatus"`
	DueDate           *time.Time `json:"dueDate"`
	DetailBillingName string     `json:"detailBillingName"`
	Amount            int64      `json:"amount"`
	BillingDetailID   uint       `json:"billingDetailId"`
}
