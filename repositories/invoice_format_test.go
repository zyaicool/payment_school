package repositories

import (
	"errors"
	"schoolPayment/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock struct
type MockInvoiceFormatRepository struct {
	mock.Mock
}

func (m *MockInvoiceFormatRepository) Create(invoiceFormat *models.InvoiceFormat) error {
	args := m.Called(invoiceFormat)
	return args.Error(0)
}

func (m *MockInvoiceFormatRepository) Update(invoiceFormat *models.InvoiceFormat) error {
	args := m.Called(invoiceFormat)
	return args.Error(0)
}

func (m *MockInvoiceFormatRepository) GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error) {
	args := m.Called(schoolID)
	return args.Get(0).(*models.InvoiceFormat), args.Error(1)
}

// Test suite struct
type InvoiceFormatRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockInvoiceFormatRepository
}

func (suite *InvoiceFormatRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(MockInvoiceFormatRepository)
}

// Test Success for Create
func (suite *InvoiceFormatRepositoryTestSuite) TestCreate_Success() {
	dummyInvoiceFormat := &models.InvoiceFormat{
		SchoolID: 1,
		Format:   "Invoice Format",
	}

	// Setting up expected behavior
	suite.mockRepo.On("Create", dummyInvoiceFormat).Return(nil)

	// Calling the repository function
	err := suite.mockRepo.Create(dummyInvoiceFormat)

	// Assertions
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for Create
func (suite *InvoiceFormatRepositoryTestSuite) TestCreate_Error() {
	dummyInvoiceFormat := &models.InvoiceFormat{
		SchoolID: 1,
		Format:   "Invoice Format",
	}

	suite.mockRepo.On("Create", dummyInvoiceFormat).Return(errors.New("create error"))

	// Calling the repository function
	err := suite.mockRepo.Create(dummyInvoiceFormat)

	// Assertions
	assert.Error(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Success for Update
func (suite *InvoiceFormatRepositoryTestSuite) TestUpdate_Success() {
	dummyInvoiceFormat := &models.InvoiceFormat{
		SchoolID: 1,
		Format:   "Updated Invoice Format",
	}

	// Setting up expected behavior
	suite.mockRepo.On("Update", dummyInvoiceFormat).Return(nil)

	// Calling the repository function
	err := suite.mockRepo.Update(dummyInvoiceFormat)

	// Assertions
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for Update
func (suite *InvoiceFormatRepositoryTestSuite) TestUpdate_Error() {
	dummyInvoiceFormat := &models.InvoiceFormat{
		SchoolID: 1,
		Format:   "Updated Invoice Format",
	}

	suite.mockRepo.On("Update", dummyInvoiceFormat).Return(errors.New("update error"))

	// Calling the repository function
	err := suite.mockRepo.Update(dummyInvoiceFormat)

	// Assertions
	assert.Error(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Success for GetBySchoolID
func (suite *InvoiceFormatRepositoryTestSuite) TestGetBySchoolID_Success() {
	schoolID := uint(1)
	dummyInvoiceFormat := &models.InvoiceFormat{
		SchoolID: schoolID,
		Format:   "Invoice Format",
	}

	// Setting up expected behavior
	suite.mockRepo.On("GetBySchoolID", schoolID).Return(dummyInvoiceFormat, nil)

	// Calling the repository function
	result, err := suite.mockRepo.GetBySchoolID(schoolID)

	// Assertions
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), dummyInvoiceFormat, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

// Test Error for GetBySchoolID
func (suite *InvoiceFormatRepositoryTestSuite) TestGetBySchoolID_Error() {
	schoolID := uint(1)

	suite.mockRepo.On("GetBySchoolID", schoolID).Return((*models.InvoiceFormat)(nil), errors.New("get error"))

	// Calling the repository function
	result, err := suite.mockRepo.GetBySchoolID(schoolID)

	// Assertions
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestInvoiceFormatRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InvoiceFormatRepositoryTestSuite))
}