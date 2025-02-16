package response

import (
	models "schoolPayment/models"
)

type TransactionListResponse struct {
	Page  int                         `json:"page"`
	Limit int                         `json:"limit"`
	Data  []models.TransactionBilling `json:"data"`
}

type MidtransResponse struct {
	OrderID     string  `json:"orderId"`
	Token       *string `json:"token"`
	RedirectURL *string `json:"redirectUrl"`
}

type MidtransExtractAPIError struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	ID            string `json:"id"`
}

type MidtransExtractResponsePayment struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	StatusCode    string   `json:"status_code,omitempty"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

type CheckingPaymentStatusResponse struct {
	Data []CheckingPaymentStatusDetailResponse `json:"data"`
}
type CheckingPaymentStatusDetailResponse struct {
	ID                   uint   `json:"id"`
	TransactionType      string `json:"transactionType"`
	VirtualAccountNumber string    `json:"virtualAccountNumber"`
	TotalAmount          int    `json:"totalAmount"`
	ReferenceNumber      string `json:"referenceNumber"`
	Description          string `json:"description"`
	StudentID            int    `json:"studentId"`
	OrderID              string `json:"orderId"`
	TransactionStatus    string `json:"transactionStatus"`
	InvoiceNumber        string `json:"invoiceNumber"`
	BillingStudentIds    string `json:"billingStudentIds"`
}

type CheckPaymentStatusFailed struct {
	ID                    uint   `json:"id"`
	TransactionType       string `json:"transactionType"`
	VirtualAccountNumber  string `json:"virtualAccountNumber"`
	TotalAmount           int    `json:"totalAmount"`
	ReferenceNumber       string `json:"referenceNumber"`
	StudentID             int    `json:"studentId"`
	OrderID               string `json:"orderId"`
	TransactionStatus     string `json:"transactionStatus"`
	InvoiceNumber         string `json:"invoiceNumber"`
	BillingStudentIds     string `json:"billingStudentIds"`
	MasterPaymentMethodID int    `json:"masterPaymentMethodId"`
	StudentName           string `json:"studentName"`
	CreatedBy             int    `json:"createdBy"`
	Nis                   string `json:"nis"`
}

type SchoolLogoSendEmailFailedResponse struct {
	SchoolName string `json:"schoolName"`
	SchoolLogo string `json:"schoolLogo"`
}
