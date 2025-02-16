package models

type Billing struct {
	Master
	BillingNumber  string       `json:"billingNumber"`
	BillingName    string       `json:"billingName"`
	BillingType    string       `json:"billingTypeId"`
	SchoolGradeID  uint         `json:"schoolGradeId"`
	SchoolYearId   uint         `json:"schoolYearId"`
	BillingAmount  int64        `json:"billingAmount"`
	Description    string       `json:"description"`
	BillingCode    string       `json:"billingCode"`
	SchoolClassIds string       `json:"schoolClassIds"`
	BankAccountId  int          `json:"bankAccountId"`
	SchoolYear     *SchoolYear  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"schoolYear"`
	SchoolGrade    *SchoolGrade `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"schoolGrade"`
	BankAccount    *BankAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bankAccount"`
	IsDonation     bool         `json:"isDonation"`
}

type BillingList struct {
	Billing                 // Embeds the original Billing struct
	CreateByUsername string `gorm:"column:create_by_username"`
}

type BillingStudentsExists struct {
	Billing
	BillingDetailID   uint   `gorm:"column:billing_detail_id"`
	DetailBillingName string `gorm:"column:detail_billing_name"`
	Amount            int64  `gorm:"column:amount"`
	IsExist           bool   `gorm:"column:is_exist"`
	Disabled          bool   `gorm:"column:disabled"`
}
