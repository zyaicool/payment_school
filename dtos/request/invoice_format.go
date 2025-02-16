package request

type CreateInvoiceFormatRequest struct {
	SchoolID               uint   `json:"schoolId" validate:"required"`
	Prefix                 string `json:"prefix" validate:"required"`
	Format                 string `json:"format" validate:"required"`
	GeneratedInvoiceFormat string `json:"generatedInvoiceFormat"`
}
