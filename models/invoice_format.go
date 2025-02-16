package models

import "time"

type InvoiceFormat struct {
	Master
	SchoolID               uint      `json:"schoolId"`
	Prefix                 string    `json:"prefix"`
	Format                 string    `json:"format"`
	GeneratedInvoiceFormat string    `json:"generatedInvoiceFormat"`
	CreatedBy              uint      `json:"createdBy"`
	UpdatedBy              uint      `json:"updatedBy"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	School                 *School   `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
}
