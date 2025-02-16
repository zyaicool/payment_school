package response

import (
	"time"
)

type BillingListResponse struct {
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	TotalPage int               `json:"totalPage"`
	TotalData int64             `json:"totalData"`
	Data      []BillingResponse `json:"data"`
}

type BillingResponse struct {
	ID              int       `json:"id"`
	BillingName     string    `json:"billingName"`
	BillingType     string    `json:"billingType"`
	BillingNumber   string    `json:"billingNumber"`
	BankAccountName string    `json:"bankAccountName"`
	SchoolGradeName string 	  `json:"schoolGradeName"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
}
