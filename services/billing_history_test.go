package services_test

import (
	"errors"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	repositories "schoolPayment/repositories"
	"schoolPayment/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repositories
type MockBillingHistoryRepository struct {
	mock.Mock
}

// Ensure the mock implements BillingHistoryRepositoryInterface
var _ repositories.BillingHistoryRepositoryInterface = (*MockBillingHistoryRepository)(nil)

type MockSchoolRepository struct {
	mock.Mock
}

type MockPaymentMethodRepository struct {
	mock.Mock
}

// MockUserBillingHistoryRepository implements the methods for mocking UserRepositoryInterface
type MockUserBillingHistoryRepository struct {
	mock.Mock
}

func (m *MockBillingHistoryRepository) GetAllBillingHistory(page int, limit int, search string, studentID int, roleID int, schoolYearId int, paymentTypeId int, schoolID int, paymentStatusCode string, sortBy string, sortOrder string, userID int, userLoginID int) ([]response.DataListBillingHistory, int, int64, error) {
	args := m.Called(page, limit, search, studentID, roleID, schoolYearId, paymentTypeId, schoolID, paymentStatusCode, sortBy, sortOrder, userID, userLoginID)
	return args.Get(0).([]response.DataListBillingHistory), args.Int(1), args.Get(2).(int64), args.Error(3)
}

func (m *MockUserBillingHistoryRepository) BulkValidateAndCreateUsers(users []models.User, userSchools []models.UserSchool) ([]models.User, []response.ResponseErrorUploadUser, error) {
	args := m.Called(users, userSchools)
	return args.Get(0).([]models.User), args.Get(1).([]response.ResponseErrorUploadUser), args.Error(2)
}

func (m *MockUserBillingHistoryRepository) DeleteAllTempVerificationEmails(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserBillingHistoryRepository) DeleteUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockBillingHistoryRepository) GetDetailBillingHistoryIDRepositories(transactionId int) (response.BillingStudentByStudentIDBillingID, error) {
	args := m.Called(transactionId)
	return args.Get(0).(response.BillingStudentByStudentIDBillingID), args.Error(1)
}

func (m *MockBillingHistoryRepository) GetInstallmentHistoryDetails(id int) (response.BillingStudentForHistory, error) {
	args := m.Called(id)
	return args.Get(0).(response.BillingStudentForHistory), args.Error(1)
}

func (m *MockBillingHistoryRepository) GetDataForInvoice(invoiceNumber string, schoolID int) ([]response.RespDataInvoice, error) {
	args := m.Called(invoiceNumber, schoolID)
	return args.Get(0).([]response.RespDataInvoice), args.Error(1)
}

func (m *MockBillingHistoryRepository) TotalAmountBillingStudent(transactionId int) (int64, error) {
	args := m.Called(transactionId)
	return args.Get(0).(int64), args.Error(1)
}

// Mock implementation of GetSchoolByID
func (m *MockSchoolRepository) GetSchoolByID(id uint) (models.School, error) {
	args := m.Called(id)
	return args.Get(0).(models.School), args.Error(1)
}

// Mock implementation of CheckNpsn
func (m *MockSchoolRepository) CheckNpsn(npsn uint) (models.School, error) {
	args := m.Called(npsn)
	return args.Get(0).(models.School), args.Error(1)
}

// Mock implementation of CreateSchool
func (m *MockSchoolRepository) CreateSchool(school *models.School) (*models.School, error) {
	args := m.Called(school)
	return args.Get(0).(*models.School), args.Error(1)
}

// Mock implementation of UpdateSchool
func (m *MockSchoolRepository) UpdateSchool(school *models.School) (*models.School, error) {
	args := m.Called(school)
	return args.Get(0).(*models.School), args.Error(1)
}

// Mock implementation of GetLastSequenceNumberSchool
func (m *MockSchoolRepository) GetLastSequenceNumberSchool() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

// Mock implementation of CheckNpsnExistsExcept
func (m *MockSchoolRepository) CheckNpsnExistsExcept(npsn uint, schoolID int) (models.School, error) {
	args := m.Called(npsn, schoolID)
	return args.Get(0).(models.School), args.Error(1)
}

// Mock implementation of GetAllSchoolList
func (m *MockSchoolRepository) GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) ([]models.SchoolList, int, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder)
	return args.Get(0).([]models.SchoolList), args.Int(1), args.Error(2)
}

// Mock implementation of GetAllOnboardingSchools
func (m *MockSchoolRepository) GetAllOnboardingSchools(search string) ([]models.School, error) {
	args := m.Called(search)
	return args.Get(0).([]models.School), args.Error(1)
}

// Mock implementation of GetSchoolsByNames
func (m *MockSchoolRepository) GetSchoolsByNames(schoolNames []string) (map[string]uint, error) {
	args := m.Called(schoolNames)
	return args.Get(0).(map[string]uint), args.Error(1)
}

// Ensure that the mock implements the SchoolRepository interface
var _ repositories.SchoolRepository = (*MockSchoolRepository)(nil)

func (m *MockPaymentMethodRepository) GetAllPaymentMethod(search string) ([]models.PaymentMethod, error) {
	args := m.Called(search)
	return args.Get(0).([]models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) GetPaymentMethodByID(id int) (*models.PaymentMethod, error) {
	args := m.Called(id)
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) CreatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	args := m.Called(paymentMethod)
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) UpdatePaymentMethod(paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	args := m.Called(paymentMethod)
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	argsCall := m.Called(query, args)
	return argsCall.Get(0).(*gorm.DB) // Pastikan mengembalikan objek yang valid
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	argsCall := m.Called(dest, conds)
	return argsCall.Get(0).(*gorm.DB) // Pastikan mengembalikan objek yang valid
}

// Tambahkan method CheckNpsn pada MockDB
func (m *MockDB) CheckNpsn(npsn uint) (models.School, error) {
	// Mendefinisikan ekspektasi mock pada method CheckNpsn
	args := m.Called(npsn)
	return args.Get(0).(models.School), args.Error(1)
}

func (m *MockUserBillingHistoryRepository) GetAllUser(page int, limit int, search string, roleID []int, schoolID int, sortBy string, sortOrder string, status *bool) ([]models.User, int, int64, error) {
	args := m.Called(page, limit, search, roleID, schoolID, sortBy, sortOrder, status)
	return args.Get(0).([]models.User), args.Int(1), args.Get(2).(int64), args.Error(3)
}

func (m *MockUserBillingHistoryRepository) GetEmailSendCount(email string, now time.Time) (int, error) {
	args := m.Called(email, now)
	return args.Int(0), args.Error(1)
}

func (m *MockUserBillingHistoryRepository) GetEmailVerification(userID int, typeFilter string) (models.TempVerificationEmail, error) {
	args := m.Called(userID, typeFilter)
	return args.Get(0).(models.TempVerificationEmail), args.Error(1)
}

func (m *MockUserBillingHistoryRepository) GetUserByIDPass(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserBillingHistoryRepository) RecordEmailSend(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserBillingHistoryRepository) UpdateTempVerificationEmail(verificationEmail *models.TempVerificationEmail) error {
	args := m.Called(verificationEmail)
	return args.Error(0)
}

func (m *MockUserBillingHistoryRepository) UpdateUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserBillingHistoryRepository) GetUserByID(id uint) (models.User, error) {
	// Buat mock data user
	user := models.User{
		Master:   models.Master{ID: id}, // Set ID melalui Master struct
		Username: "mockuser",
		Email:    "mockuser@example.com", // Misalnya, menambahkan email
	}
	return user, nil
}

func (m *MockDB) CheckNpsnExistsExcept(npsn uint, schoolID int) (models.School, error) {
	args := m.Called(npsn, schoolID)
	return args.Get(0).(models.School), args.Error(1)
}

func (m *MockDB) CreateSchool(school *models.School) (*models.School, error) {
	args := m.Called(school)
	return args.Get(0).(*models.School), args.Error(1)
}

func (m *MockDB) GetAllOnboardingSchools(schoolID string) ([]models.School, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]models.School), args.Error(1)
}

func (m *MockDB) GetAllSchoolList(page int, limit int, search string, sortBy string, sortOrder string) ([]models.SchoolList, int, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder)
	return args.Get(0).([]models.SchoolList), args.Int(1), args.Error(2)
}

func (m *MockDB) GetLastSequenceNumberSchool() (int, error) {
	args := m.Called() // No argument passed here
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetSchoolByID(id uint) (models.School, error) {
	args := m.Called(id)
	return args.Get(0).(models.School), args.Error(1)
}

func (m *MockDB) GetSchoolsByNames(schoolNames []string) (map[string]uint, error) {
	args := m.Called(schoolNames)
	return args.Get(0).(map[string]uint), args.Error(1)
}

func (m *MockDB) UpdateSchool(school *models.School) (*models.School, error) {
	args := m.Called(school)
	return args.Get(0).(*models.School), args.Error(1)
}

func TestGetAllBillingHistory(t *testing.T) {
	// Create mock repositories
	mockBillingHistoryRepo := new(MockBillingHistoryRepository)
	mockUserRepo := new(MockUserBillingHistoryRepository)
	mockSchoolRepo := new(MockSchoolRepository)
	mockPaymentRepo := new(MockPaymentMethodRepository)
	mockDB := new(MockDB)

	// Mocking the DB query for service
	mockDB.On("Table", "billing_histories").Return(mockDB)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(mockDB)
	mockDB.On("First", mock.Anything, mock.Anything).Return(mockDB)

	// Mocking the GetUserByID call
	t.Log("Mocking GetUserByID call")
	mockUserRepo.On("GetUserByID", uint(1)).Return(models.User{
		Master:   models.Master{ID: 1}, // Use non-pointer if not needed
		Username: "john.doe",
		Email:    "john.doe@example.com",
		UserSchool: &models.UserSchool{
			SchoolID: 1,
		},
		IsVerification: false,
		IsBlock:        false,
		Image:          "image.jpg",
	}, nil)
	t.Log("Mocked GetUserByID successfully")

	// Mocking GetAllBillingHistory call
	now := time.Now()
	t.Log("Mocking GetAllBillingHistory response")
	mockBillingHistoryRepo.On("GetAllBillingHistory", 1, 10, "", 0, 1, 0, 0, 0, "", "", "", 1, 1).Return(
		[]response.DataListBillingHistory{
			{
				ID:                1,
				InvoiceNumber:     "INV123",
				StudentName:       "John Doe",
				PaymentDate:       &now, // Ensure it's a valid time pointer
				PaymentMethod:     "kasir",
				Username:          "admin",
				TotalAmount:       1000,
				TransactionStatus: "menunggu",
				OrderID:           "order1",
			},
		},
		1,
		int64(100),
		nil,
	)
	t.Log("Mocked GetAllBillingHistory successfully")

	// Mocking the CheckNpsn method
	mockSchoolRepo.On("CheckNpsn", uint(1)).Return(models.School{}, errors.New("School not found"))

	// Create the service instance, passing in all required mocks
	service := services.NewBillingHistoryService(
		mockBillingHistoryRepo,
		mockUserRepo,
		mockSchoolRepo,  // Passing the mockSchoolRepo to the service
		mockPaymentRepo, // Passing the mockPaymentRepo to the service
	)

	// // Check if the service is initialized correctly
	// t.Log("Checking service initialization")
	// if service != nil {
	// 	t.Fatal("Service is nil!")
	// }

	// Call the method
	// Anda dapat memanggil method untuk melihat apakah objek berfungsi:
	result, err := service.GetAllBillingHistory(
		1,  // page
		10, // limit
		"", // search
		0,  // studentID
		1,  // roleID
		0,  // schoolYearId
		0,  // paymentTypeId
		"", // paymentStatusCode
		1,  // userID
		1,  // userLoginID
		"", // sortBy
		"", // sortOrder
	)

	t.Log("Service result:", result) // Bisa digunakan untuk verifikasi jika service berfungsi
	assert.NoError(t, err)

	// Debugging output
	t.Log("Service result:", result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, int64(100), result.TotalData)
	assert.Equal(t, 1, result.TotalPage)
	assert.Len(t, result.Data, 1)

	// Verify all mock expectations
	mockUserRepo.AssertExpectations(t)
	mockBillingHistoryRepo.AssertExpectations(t)
	mockSchoolRepo.AssertExpectations(t)
	mockPaymentRepo.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
