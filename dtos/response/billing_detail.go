package response

import "time"

type BillingDetailResponse struct {
	ID              uint            `json:"id"`              // ID of the billing
	BillingName     string          `json:"billingName"`     // Name of the billing
	BillingCode     string          `json:"billingCode"`     // Code of the billing
	BillingType     string          `json:"billingType"`     // Type of the billing
	BankAccountName string          `json:"bankAccountName"` // Bank account associated with the billing
	Description     string          `json:"description"`     // Description of the billing
	SchoolYear      string          `json:"schoolYear"`      // Name of the school year
	SchoolClassList string          `json:"schoolClassList"` // List of school classes (if applicable)
	DetailBillings  []DetailBilling `json:"detailBillings"`
}

type DetailBilling struct {
	ID                uint      `json:"id"` // Add ID field
	DetailBillingName string    `json:"detailBillingName"`
	DueDate           time.Time `json:"dueDate"`
	Amount            int64     `json:"amount"`
}
