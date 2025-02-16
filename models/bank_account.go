package models

type BankAccount struct {
	Master
	SchoolID      uint    `json:"schoolId"`
	BankName      string  `json:"bankName"`
	AccountName   string  `json:"accountName"`
	AccountNumber string  `json:"accountNumber"`
	AccountOwner  string  `json:"accountOwner"`
	School        *School `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
}
