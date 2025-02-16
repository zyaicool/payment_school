package services

import (
	"errors"
	request "schoolPayment/dtos/request"
	"schoolPayment/models"
	"schoolPayment/repositories"
	"time"

	"gorm.io/gorm"
)


type InvoiceFormatServiceInterface interface {
	Create(request *request.CreateInvoiceFormatRequest, userID int) (*models.InvoiceFormat, error)
	GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error)
}

type InvoiceFormatService struct {
	invoiceFormatRepo repositories.InvoiceFormatRepositoryInterface
}

var _ InvoiceFormatServiceInterface = (*InvoiceFormatService)(nil)

func NewInvoiceFormatService(invoiceFormatRepo repositories.InvoiceFormatRepositoryInterface) *InvoiceFormatService {
	return &InvoiceFormatService{
		invoiceFormatRepo: invoiceFormatRepo,
	}
}


func (s *InvoiceFormatService) Create(request *request.CreateInvoiceFormatRequest, userID int) (*models.InvoiceFormat, error) {
	
	existingFormat, err := s.invoiceFormatRepo.GetBySchoolID(request.SchoolID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	
	now := time.Now()

	
	if existingFormat != nil {
		existingFormat.Prefix = request.Prefix
		existingFormat.Format = request.Format
		existingFormat.GeneratedInvoiceFormat = request.GeneratedInvoiceFormat
		existingFormat.UpdatedAt = now
		existingFormat.UpdatedBy = uint(userID)

		err = s.invoiceFormatRepo.Update(existingFormat)
		if err != nil {
			return nil, err
		}
		return existingFormat, nil
	}

	
	invoiceFormat := &models.InvoiceFormat{
		SchoolID:               request.SchoolID,
		Prefix:                 request.Prefix,
		Format:                 request.Format,
		GeneratedInvoiceFormat: request.GeneratedInvoiceFormat,
		CreatedAt:              now,
		CreatedBy:              uint(userID),
		UpdatedAt:              now,
		UpdatedBy:              uint(userID),
	}

	err = s.invoiceFormatRepo.Create(invoiceFormat)
	if err != nil {
		return nil, err
	}

	return invoiceFormat, nil
}

func (s *InvoiceFormatService) GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error) {
	return s.invoiceFormatRepo.GetBySchoolID(schoolID)
}
