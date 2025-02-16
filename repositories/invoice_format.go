package repositories

import (
	"schoolPayment/models"

	"gorm.io/gorm"
)

type InvoiceFormatRepositoryInterface interface {
	GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error)
	Update(invoiceFormat *models.InvoiceFormat) error
	Create(invoiceFormat *models.InvoiceFormat) error
}

type InvoiceFormatRepository struct {
	db *gorm.DB
}

func NewInvoiceFormatRepository(db *gorm.DB) *InvoiceFormatRepository {
	return &InvoiceFormatRepository{db: db}
}

func (r *InvoiceFormatRepository) Create(invoiceFormat *models.InvoiceFormat) error {
	return r.db.Create(invoiceFormat).Error
}

func (r *InvoiceFormatRepository) Update(invoiceFormat *models.InvoiceFormat) error {
	return r.db.Save(invoiceFormat).Error
}

func (r *InvoiceFormatRepository) GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error) {
	var invoiceFormat models.InvoiceFormat
	err := r.db.Where("school_id = ?", schoolID).First(&invoiceFormat).Error
	if err != nil {
		return nil, err
	}
	return &invoiceFormat, nil
}
