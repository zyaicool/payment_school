package request

type MidtransRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CreditCard         CreditCard         `json:"credit_card"`
}

type TransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type CreditCard struct {
	Secure bool `json:"secure"`
}

type CreateTransactionRequest struct {
	BillingId         int      `json:"billingId"`
	StudentId         int      `json:"studentId"`
	AmountToPay       int      `json:"amountToPay"`
	BillingStudentIds []string `json:"billingStudentIds"`
	Discount          int      `json:"discount"`
	DiscountType      string   `json:"discountType"`
	ChangeAmount      int      `json:"changeAmount"`
	TotalBilling      int      `json:"totalBilling"`
	Description       string   `json:"description"`
	PaymentMethodId   int      `json:"paymentMethodId"`
	Amount            int      `json:"amount"`
}

type SubmitTransactionRequest struct {
	MerchantID        string `json:"merchantId"`
	OrderID           string `json:"orderId"`
	TransactionStatus string `json:"transactionStatus"`
	StatusCode        string `json:"statusCode"`
}
