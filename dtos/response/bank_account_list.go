package response

import (
	"time"
)

type BankAccountListResponse struct {
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	TotalPage int               `json:"totalPage"` // Added totalPage field
	TotalData int               `json:"totalData"` // Added totalData field
	Data      []BankAccountData `json:"data"`
}

type BankAccountData struct {
	ID            int        `json:"id"`
	BankName      string     `json:"bankName"`
	AccountNumber string     `json:"accountNumber"`
	AccountName   string     `json:"accountName"`
	AccountOwner  string     `json:"accountOwner"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     int        `json:"updatedBy"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	School        SchoolData `json:"school"`      // Include the school data
	PlaceHolder   string     `json:"placeholder"` // Add placeholder field
	IsDelete      bool       `json:"isDelete"`
}

type SchoolData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Add other fields as necessary
}
