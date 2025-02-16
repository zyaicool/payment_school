package request

type BillingTypeCreateUpdateRequest struct {
	SchoolID          int    `json:"schoolId"`
	BillingTypeName   string `json:"billingTypeName"`
	BillingTypePeriod string `json:"billingTypePeriod"`
}
