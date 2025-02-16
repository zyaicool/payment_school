package request

type WebhookPayload struct {
	VANumbers         []VANumber `json:"va_numbers"`
	TransactionTime   string     `json:"transaction_time"`
	TransactionStatus string     `json:"transaction_status"`
	TransactionID     string     `json:"transaction_id"`
	StatusMessage     string     `json:"status_message"`
	StatusCode        string     `json:"status_code"`
	SignatureKey      string     `json:"signature_key"`
	SettlementTime    *string    `json:"settlement_time"` // Pointer to make it optional
	PaymentType       string     `json:"payment_type"`
	PaymentAmounts    []string   `json:"payment_amounts"`
	OrderID           string     `json:"order_id"`
	MerchantID        string     `json:"merchant_id"`
	GrossAmount       string     `json:"gross_amount"`
	FraudStatus       string     `json:"fraud_status"`
	ExpiryTime        string     `json:"expiry_time"`
	Currency          string     `json:"currency"`
	PermataVaNumber   string     `json:"permata_va_number"`
	BillerCode        string     `json:"biller_code"`
	BillKey           string     `json:"bill_key"`
}

type VANumber struct {
	VANumber string `json:"va_number"`
	Bank     string `json:"bank"`
}
