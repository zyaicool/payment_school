package request

type BankAccountCreateRequest struct {
	SchoolID      uint   `json:"schoolId"`
	BankName      string `json:"bankName"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	AccountOwner  string `json:"accountOwner"`
}
