package models

type BillingType struct {
	Master
	SchoolID          uint    `json:"schoolId"`
	BillingTypeCode   string  `json:"billingTypeCode"`
	BillingTypeName   string  `json:"billingTypeName"`
	BillingTypePeriod string  `json:"billingTypePeriod"`
	School            *School `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
}
