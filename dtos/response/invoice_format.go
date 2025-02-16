package response

import "time"

type InvoiceFormatResponse struct {
	ID                     uint      `json:"id"`
	SchoolID               uint      `json:"schoolId"`
	Prefix                 string    `json:"prefix"`
	Format                 string    `json:"format"`
	GeneratedInvoiceFormat string    `json:"generatedInvoiceFormat"`
	CreatedBy              uint      `json:"createdBy"`
	UpdatedBy              uint      `json:"updatedBy"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}
