package models

type TransactionBillingHistory struct {
	Master
	TransactionBillingId uint   `json:"transactionBillingId"`
	ReferenceNumber      string `json:"referenceNumber"`
	OrderID              string `json:"orderId"`
	InvoiceNumber        string `json:"invoiceNumber"`
	TransactionStatus    string `json:"transactionStatus"`
}
