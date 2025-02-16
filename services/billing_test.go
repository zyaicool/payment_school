package services

import (
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)
type MockBillingService struct {
	mock.Mock
}

func (m *MockBillingService) GetBillingByID(id uint) (response.BillingDetailResponse, error) {
	args := m.Called(id)
	return args.Get(0).(response.BillingDetailResponse), args.Error(1)
}

type MockBillingRepository struct {
	mock.Mock
}

func (m *MockBillingRepository) GetDetailBillingsByBillingID(billingID uint) ([]response.DetailBilling, error) {
	args := m.Called(billingID)
	return args.Get(0).([]response.DetailBilling), args.Error(1)
}

func (m *MockBillingRepository) GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, user models.User) ([]models.BillingList, int, int64, error) {
	args := m.Called(page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, isDonation, user)
	return args.Get(0).([]models.BillingList), args.Int(1), args.Get(2).(int64), args.Error(3)
}

func (m *MockBillingRepository) GetBillingByID(id int) (models.Billing, error) {
	args := m.Called(id)
	return args.Get(0).(models.Billing), args.Error(1)
}

func (m *MockBillingRepository) CheckBillingCode(billingCode string) bool {
	args := m.Called(billingCode)
	return args.Bool(0)
}

func (m *MockBillingRepository) CreateBillingDonation(billing *models.Billing) (*models.Billing, error) {
	args := m.Called(billing)
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingRepository) GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) ([]models.BillingStudentsExists, error) {
	args := m.Called(studentID, schoolYearID, schoolGradeID, schoolClassID)
	return args.Get(0).([]models.BillingStudentsExists), args.Error(1)
}

func (m *MockBillingRepository) CheckBillingStudentExists(studentID uint, billingDetailID uint) (bool, error) {
	args := m.Called(studentID, billingDetailID)
	return args.Bool(0), args.Error(1)
}

func (m *MockBillingRepository) CreateBilling(billing *models.Billing) (*models.Billing, error) {
	args := m.Called(billing)
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingRepository) GetLastSequenceNumberBilling() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func TestGetBillingByID(t *testing.T) {
	// Prepare the mock service
	mockService := new(MockBillingService)

	// Set up mock response
	expectedResponse := response.BillingDetailResponse{
		ID:              1,
		BillingName:     "Tuition Fee",
		BillingCode:     "TF123",
		BillingType:     "Regular",
		BankAccountName: "Bank XYZ - 123456",
		Description:     "Tuition fee for the semester",
		SchoolYear:      "2024/2025",
		SchoolClassList: "Class A, Class B",
	}

	// Set up expectation
	mockService.On("GetBillingByID", uint(1)).Return(expectedResponse, nil)

	// Call the method
	result, err := mockService.GetBillingByID(1)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, result)

	// Ensure that the method was called with the correct parameters
	mockService.AssertExpectations(t)
}

func TestCreateBilling_Success(t *testing.T) {
	mockSchoolGradeRepo := new(MockSchoolGradeRepository)
	mockSchoolYearRepo := new(MockSchoolYearRepository)
	mockUserRepo := new(MockUserRepository)
	mockBillingRepo := new(MockBillingRepository)
	mockStudentRepo := new(MockStudentRepository)

	mockSchoolGradeRepo.On("GetSchoolGradeByID", uint(1)).Return(&models.SchoolGrade{SchoolGradeCode: "SMK"}, nil).Once()
	mockSchoolYearRepo.On("GetSchoolYearByID", uint(2)).Return(&models.SchoolYear{SchoolYearName: "2023/2024"}, nil).Once()
	mockUserRepo.On("GetUserByID", uint(1)).Return(&models.User{UserSchool: &models.UserSchool{ SchoolID: 1 }}, nil).Once()
	mockBillingRepo.On("GetLastSequenceNumberBilling").Return(1, nil).Once()
	mockBillingRepo.On("CheckBillingCode", "BILL123").Return(true).Once()
	mockBillingRepo.On("CreateBilling", mock.Anything).Return(&models.Billing{
		BillingName:    "Test Billing",
		BillingCode:    "BILL123",
		BillingType:    "Monthly",
		SchoolGradeID:  1,
		SchoolYearId:   2,
		BillingAmount:  1000,
		Description:    "Test billing description",
		BankAccountId:  1,
		SchoolClassIds: "1,2",
	}, nil).Once()

	// Mock GetAllStudentForBilling
	mockStudents := []models.Student{
		{Nisn: "1234"},
		{Nisn: "1235"},
	}
	mockUser := models.User{UserSchool: &models.UserSchool{SchoolID: 1}}
	schoolGradeID := 1
	schoolClassIDs := []int{1, 2}

	mockStudentRepo.On("GetAllStudentForBilling", mockUser, schoolGradeID, schoolClassIDs).
		Return(mockStudents, nil).Once()

	billingService := BillingService{
		billingRepository: mockBillingRepo,
		userRepository: mockUserRepo,
		schoolYearRepository: mockSchoolYearRepo,
		schoolGradeRepository: mockSchoolGradeRepo,
		studentRepository: mockStudentRepo,
	}

	// Define request
	request := &request.BillingCreateRequest{
		BillingName:    "Test Billing",
		BillingCode:    "BILL123",
		BillingType:    "Monthly",
		SchoolGradeID:  1,
		SchoolYearId:   2,
		BillingAmount:  1000,
		Description:    "Test billing description",
		BankAccountId:  1,
		SchoolClassIds: []string{"1", "2"},
		DetailBillings: []request.DetailBillings{
			{Amount: 500},
			{Amount: 500},
		},
	}

	// Call service method
	result, err := billingService.CreateBilling(request, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Billing", result.BillingName)
	assert.Equal(t, "BILL123", result.BillingCode)
	assert.Equal(t, "Monthly", result.BillingType)

	// Verify mock expectations
	mockSchoolYearRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockBillingRepo.AssertExpectations(t)
}

func TestGetAllBilling_Success(t *testing.T) {
	// Initialize mocks
	mockUserRepo := new(MockUserRepository)
	mockBillingRepo := new(MockBillingRepository)

	// Create the service with mocked repositories
	billingService := &BillingService{
		userRepository:    mockUserRepo,
		billingRepository: mockBillingRepo,
	}

	// Input parameters
	page := 1
	limit := 10
	search := "Test"
	billingType := "Type1"
	paymentType := "Online"
	schoolGrade := "1"
	sort := "ASC"
	sortBy := "createdAt"
	sortOrder := "ASC"
	bankAccountId := 1
	isDonation := new(bool)
	*isDonation = false
	userID := 1

	// Mocked user data
	mockUser := models.User{
		Master: models.Master{
			ID: uint(userID),
		},
		UserSchool: &models.UserSchool{
			SchoolID: 1,
		},
	}
	mockUserRepo.On("GetUserByID", uint(userID)).Return(mockUser, nil)

	// Mocked billing data
	mockBillingList := []models.BillingList{
		{
			Billing: models.Billing{
				Master: models.Master{
					ID: 1,
				},
				BillingName:  "Test Billing",
				BillingType:  "Type1",
				BankAccount:  &models.BankAccount{BankName: "TestBank", AccountNumber: "123456789"},
				SchoolGrade:  &models.SchoolGrade{SchoolGradeName: "Grade 1"},
			},
		},
	}
	mockTotalPages := 10
	mockTotalData := int64(100)

	mockBillingRepo.On("GetAllBilling", page, limit, search, billingType, paymentType, schoolGrade, sort, "created_at", sortOrder, bankAccountId, isDonation, mock.Anything).
		Return(mockBillingList, mockTotalPages, mockTotalData, nil)

	// Call the service method
	result, err := billingService.GetAllBilling(page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, isDonation, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, page, result.Page)
	assert.Equal(t, mockTotalData, int64(result.TotalData))
	assert.Equal(t, mockTotalPages, result.TotalPage)
	assert.Len(t, result.Data, 1)

	// Validate the response data
	billingResponse := result.Data[0]
	assert.Equal(t, 1, billingResponse.ID)
	assert.Equal(t, "Test Billing", billingResponse.BillingName)
	assert.Equal(t, "Type1", billingResponse.BillingType)
	assert.Equal(t, "TestBank - 123456789", billingResponse.BankAccountName)
	assert.Equal(t, "Grade 1", billingResponse.SchoolGradeName)

	// Ensure all mocks are called
	mockUserRepo.AssertExpectations(t)
	mockBillingRepo.AssertExpectations(t)
}
