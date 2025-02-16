package models

type TransactionBilling struct {
	Master
	BillingID            string                      `json:"billingId"`
	StudentID            uint                        `json:"studentId"`
	TransactionType      string                      `json:"transactionType"`
	VirtualAccountNumber string                      `json:"virtualAccountNumber"`
	TotalAmount          int                         `json:"totalAmount"`
	ReferenceNumber      string                      `json:"referenceNumber"`
	Description          string                      `json:"description"`
	OrderID              string                      `json:"orderId"`
	TransactionStatus    string                      `json:"transactionStatus"`
	InvoiceNumber        string                      `json:"invoiceNumber"`
	BillingStudentIds    string                      `json:"billingStudentIds"`
	AccountNumber        string                      `json:"accountNumber"`
	ExpiryTime           string                      `json:"expiryTime"`
	TransactionHistory   []TransactionBillingHistory `gorm:"foreignKey:TransactionBillingId"`
}
