package services_test

import (
	"schoolPayment/dtos/request"
	"schoolPayment/models"
	"schoolPayment/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockInvoiceFormatRepository struct {
	mock.Mock
}

func (m *MockInvoiceFormatRepository) GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error) {
	args := m.Called(schoolID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InvoiceFormat), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceFormatRepository) Update(invoiceFormat *models.InvoiceFormat) error {
	args := m.Called(invoiceFormat)
	return args.Error(0)
}

func (m *MockInvoiceFormatRepository) Create(invoiceFormat *models.InvoiceFormat) error {
	args := m.Called(invoiceFormat)
	return args.Error(0)
}

func TestInvoiceFormatService_Create(t *testing.T) {
	mockRepo := new(MockInvoiceFormatRepository)
	service := services.NewInvoiceFormatService(mockRepo)

	// Test case: Update existing format
	existingFormat := &models.InvoiceFormat{
		SchoolID:               1,
		Prefix:                 "INV",
		Format:                 "YYYY-MM",
		GeneratedInvoiceFormat: "INV-2024-12",
	}
	mockRepo.On("GetBySchoolID", uint(1)).Return(existingFormat, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	requestInvoiceFormat:= &request.CreateInvoiceFormatRequest{
		SchoolID:               1,
		Prefix:                 "NEW",
		Format:                 "MM-YYYY",
		GeneratedInvoiceFormat: "NEW-12-2024",
	}
	
	updatedFormat, err := service.Create(requestInvoiceFormat, 123)

	assert.NoError(t, err)
	assert.Equal(t, "NEW", updatedFormat.Prefix)
	assert.Equal(t, "MM-YYYY", updatedFormat.Format)
	assert.Equal(t, "NEW-12-2024", updatedFormat.GeneratedInvoiceFormat)
	mockRepo.AssertCalled(t, "Update", mock.Anything)

	// Test case: Create new format
	mockRepo.On("GetBySchoolID", uint(2)).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.Anything).Return(nil)

	newRequest := &request.CreateInvoiceFormatRequest{
		SchoolID:               2,
		Prefix:                 "NEW",
		Format:                 "DD-MM-YYYY",
		GeneratedInvoiceFormat: "NEW-30-12-2024",
	}

	newFormat, err := service.Create(newRequest, 123)

	assert.NoError(t, err)
	assert.Equal(t, uint(2), newFormat.SchoolID)
	assert.Equal(t, "NEW", newFormat.Prefix)
	assert.Equal(t, "DD-MM-YYYY", newFormat.Format)
	assert.Equal(t, "NEW-30-12-2024", newFormat.GeneratedInvoiceFormat)
	mockRepo.AssertCalled(t, "Create", mock.Anything)

	// Test case: Error retrieving format
	mockRepo.On("GetBySchoolID", uint(3)).Return(nil, gorm.ErrInvalidTransaction)

	errRequest := &request.CreateInvoiceFormatRequest{
		SchoolID:               3,
		Prefix:                 "ERR",
		Format:                 "ERR-FORMAT",
		GeneratedInvoiceFormat: "ERR-123",
	}

	_, err = service.Create(errRequest, 123)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	mockRepo.AssertCalled(t, "GetBySchoolID", uint(3))
}
