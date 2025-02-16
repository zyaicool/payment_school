package response

import (
	"math/big"
	"time"
)

type PaymentReportResponse struct {
	TotalTransactionAmount *big.Int          `json:"totalTransactionAmount"`
	TotalTransaction       int64             `json:"totalTransaction"`
	TotalStudent           int               `json:"totalStudent"`
	ListPaymentReport      ListPaymentReport `json:"listPaymentReport"`
}

type ListPaymentReport struct {
	Page      int                   `json:"page"`
	Limit     int                   `json:"limit"`
	TotalPage int                   `json:"totalPage"`
	TotalData int64                 `json:"totalData"`
	Data      []PaymentReportDetail `json:"data"`
}

type PaymentReportDetail struct {
	ID                int        `json:"id"`
	InvoiceNumber     string     `json:"invoiceNumber"`
	StudentName       string     `json:"studentName"`
	PaymentDate       *time.Time `json:"paymentDate"`
	PaymentMethod     string     `json:"paymentMethod"`
	Username          string     `json:"username"`
	SchoolGradeName   string     `json:"schoolGradeName"`
	SchoolClassName   string     `json:"schoolClassName"`
	TotalAmount       int64      `json:"totalAmount"`
	TransactionStatus string     `json:"transactionStatus"`
}

// New struct for SQL query scanning
type PaymentReportSummary struct {
	TotalTransactionAmount string `json:"totalTransactionAmount"`
	TotalTransaction       int64  `json:"totalTransaction"`
	TotalStudent           int    `json:"totalStudent"`
}
