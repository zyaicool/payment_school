package services

import (
	"errors"
	"schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolClassRepository struct {
	mock.Mock
}

func (m *MockSchoolClassRepository) CreateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	args := m.Called(schoolClass)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassRepository) GetLastSequenceNumberSchoolClasss() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSchoolClassRepository) GetSchoolClassByID(id uint) (models.SchoolClass, error) {
	args := m.Called(id)
	return args.Get(0).(models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassRepository) UpdateSchoolClass(schoolClass *models.SchoolClass) (*models.SchoolClass, error) {
	args := m.Called(schoolClass)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

// GetUserByID mocks the repository method for getting a user by ID
func (m *MockUserRepository) GetUserByID(id uint) (models.User, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*models.User); ok {
		return *user, args.Error(1)
	}
	return models.User{}, args.Error(1)
}

// BulkValidateAndCreateUsers mocks the repository method for bulk validating and creating users
func (m *MockUserRepository) BulkValidateAndCreateUsers(users []models.User, userSchools []models.UserSchool) ([]models.User, []response.ResponseErrorUploadUser, error) {
	args := m.Called(users, userSchools)
	return args.Get(0).([]models.User), args.Get(1).([]response.ResponseErrorUploadUser), args.Error(2)
}

// DeleteAllTempVerificationEmails mocks the repository method for deleting all temporary verification emails
func (m *MockUserRepository) DeleteAllTempVerificationEmails(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

// DeleteUserRepository mocks the repository method for deleting a user
func (m *MockUserRepository) DeleteUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// GetAllUser mocks the repository method for getting all users
func (m *MockUserRepository) GetAllUser(page int, limit int, search string, roleIDs []int, schoolID int, sortBy string, sortOrder string, status *bool) ([]models.User, int, int64, error) {
	args := m.Called(page, limit, search, roleIDs, schoolID, sortBy, sortOrder, status)
	return args.Get(0).([]models.User), args.Get(1).(int), args.Get(2).(int64), args.Error(3)
}

// GetEmailSendCount mocks the repository method for getting email send count
func (m *MockUserRepository) GetEmailSendCount(email string, startTime time.Time) (int, error) {
	args := m.Called(email, startTime)
	return args.Get(0).(int), args.Error(1)
}

// GetEmailVerification mocks the repository method for getting email verification
func (m *MockUserRepository) GetEmailVerification(userID int, email string) (models.TempVerificationEmail, error) {
	args := m.Called(userID, email)
	return args.Get(0).(models.TempVerificationEmail), args.Error(1)
}

// GetUserByEmail mocks the repository method for getting a user by email
func (m *MockUserRepository) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

// CreateUser mocks the repository method for creating a user
func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// UpdateUser mocks the repository method for updating a user
func (m *MockUserRepository) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// CreateTempVerificationEmail mocks the repository method for creating a temporary verification email
func (m *MockUserRepository) CreateTempVerificationEmail(tempVerificationEmail *models.TempVerificationEmail) (*models.TempVerificationEmail, error) {
	args := m.Called(tempVerificationEmail)
	return args.Get(0).(*models.TempVerificationEmail), args.Error(1)
}

// UpdateTempVerificationEmail mocks the repository method for updating a temporary verification email
func (m *MockUserRepository) UpdateTempVerificationEmail(tempVerificationEmail *models.TempVerificationEmail) error {
	args := m.Called(tempVerificationEmail)
	return args.Error(0)
}

// GetUserByIDPass mocks the repository method for getting a user by ID pass
func (m *MockUserRepository) GetUserByIDPass(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

// RecordEmailSend mocks the repository method for recording email send
func (m *MockUserRepository) RecordEmailSend(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

// UpdateUserRepository mocks the repository method for updating a user
func (m *MockUserRepository) UpdateUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateSchoolClass(t *testing.T) {
	mockRepo := new(MockSchoolClassRepository)
	mockUserRepo := new(MockUserRepository)

	schoolClassService := NewSchoolClassService(mockRepo, mockUserRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetLastSequenceNumberSchoolClasss").Return(10, nil).Once()
		mockRepo.On("CreateSchoolClass", mock.Anything).Return(&models.SchoolClass{
			SchoolClassName: "Class 1",
			SchoolClassCode: "SC011",
		}, nil).Once()

		schoolClassRequest := &request.SchoolClassCreateUpdateRequest{
			SchoolID:        1,
			SchoolGradeID:   2,
			SchoolClassName: "Class 1",
			Suffix:          "A",
			SchoolMajorID:   3,
			PrefixClassID:   4,
		}

		result, err := schoolClassService.CreateSchoolClass(schoolClassRequest, 123)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "SC011", result.SchoolClassCode)
		assert.Equal(t, "Class 1", result.SchoolClassName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Generating Code", func(t *testing.T) {
		mockRepo.On("GetLastSequenceNumberSchoolClasss").Return(0, errors.New("error fetching sequence")).Once()

		schoolClassRequest := &request.SchoolClassCreateUpdateRequest{
			SchoolID:        1,
			SchoolGradeID:   2,
			SchoolClassName: "Class 1",
			Suffix:          "A",
			SchoolMajorID:   3,
			PrefixClassID:   4,
		}

		result, err := schoolClassService.CreateSchoolClass(schoolClassRequest, 123)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// MockUserClassRepository implements the methods for mocking UserRepositoryInterface
type MockUserClassRepository struct {
	mock.Mock
}

// Definisikan ulang metode-metode mock untuk MockUserClassRepository
func (m *MockUserClassRepository) GetUserByID(id uint) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

// Tambahkan implementasi metode BulkValidateAndCreateUsers
func (m *MockUserClassRepository) BulkValidateAndCreateUsers(users []models.User, userSchools []models.UserSchool) ([]models.User, []response.ResponseErrorUploadUser, error) {
	args := m.Called(users, userSchools)
	return args.Get(0).([]models.User), args.Get(1).([]response.ResponseErrorUploadUser), args.Error(2)
}

// DeleteAllTempVerificationEmails mocks the repository method for deleting all temporary verification emails
func (m *MockUserClassRepository) DeleteAllTempVerificationEmails(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

// DeleteUserRepository mocks the repository method for deleting a user
func (m *MockUserClassRepository) DeleteUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// GetAllUser mocks the repository method for getting all users
func (m *MockUserClassRepository) GetAllUser(page int, limit int, search string, roleIDs []int, schoolID int, sortBy string, sortOrder string, status *bool) ([]models.User, int, int64, error) {
	args := m.Called(page, limit, search, roleIDs, schoolID, sortBy, sortOrder, status)
	return args.Get(0).([]models.User), args.Get(1).(int), args.Get(2).(int64), args.Error(3)
}

// GetEmailSendCount mocks the repository method for getting email send count
func (m *MockUserClassRepository) GetEmailSendCount(email string, startTime time.Time) (int, error) {
	args := m.Called(email, startTime)
	return args.Get(0).(int), args.Error(1)
}

// GetEmailVerification mocks the repository method for getting email verification
func (m *MockUserClassRepository) GetEmailVerification(userID int, email string) (models.TempVerificationEmail, error) {
	args := m.Called(userID, email)
	return args.Get(0).(models.TempVerificationEmail), args.Error(1)
}

// GetUserByEmail mocks the repository method for getting a user by email
func (m *MockUserClassRepository) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

// CreateUser mocks the repository method for creating a user
func (m *MockUserClassRepository) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// UpdateUser mocks the repository method for updating a user
func (m *MockUserClassRepository) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// CreateTempVerificationEmail mocks the repository method for creating a temporary verification email
func (m *MockUserClassRepository) CreateTempVerificationEmail(tempVerificationEmail *models.TempVerificationEmail) (*models.TempVerificationEmail, error) {
	args := m.Called(tempVerificationEmail)
	return args.Get(0).(*models.TempVerificationEmail), args.Error(1)
}

// UpdateTempVerificationEmail mocks the repository method for updating a temporary verification email
func (m *MockUserClassRepository) UpdateTempVerificationEmail(tempVerificationEmail *models.TempVerificationEmail) error {
	args := m.Called(tempVerificationEmail)
	return args.Error(0)
}

// GetUserByIDPass mocks the repository method for getting a user by ID pass
func (m *MockUserClassRepository) GetUserByIDPass(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

// RecordEmailSend mocks the repository method for recording email send
func (m *MockUserClassRepository) RecordEmailSend(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

// UpdateUserRepository mocks the repository method for updating a user
func (m *MockUserClassRepository) UpdateUserRepository(user *models.User) (*models.User, error) {
	args := m.Called(user)
	return args.Get(0).(*models.User), args.Error(1)
}

// TestDeleteSchoolClass_Success tests the successful deletion of a school class
func TestDeleteSchoolClass_Success(t *testing.T) {
	// Create mock repositories
	mockSchoolClassRepo := new(MockSchoolClassRepository)
	mockUserSchoolClassRepo := new(MockUserClassRepository) // mock untuk repository User

	// Initialize the service with mock repositories
	schoolClassService := NewSchoolClassService(mockSchoolClassRepo, mockUserSchoolClassRepo)

	// Setup test data
	userID := 1
	classID := uint(1)
	mockSchoolClass := models.SchoolClass{
		Master: models.Master{ID: classID},
	}

	// Expected data after deletion
	updatedSchoolClass := mockSchoolClass
	currentTime := time.Now()
	updatedSchoolClass.Master.DeletedAt = &currentTime
	updatedSchoolClass.Master.DeletedBy = &userID

	// Define mock behavior
	mockSchoolClassRepo.On("GetSchoolClassByID", classID).Return(mockSchoolClass, nil)
	mockSchoolClassRepo.On("UpdateSchoolClass", &updatedSchoolClass).Return(&updatedSchoolClass, nil)

	// Call the method under test
	deletedClass, err := schoolClassService.DeleteSchoolClass(classID, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, deletedClass)
	assert.Equal(t, updatedSchoolClass.Master.DeletedAt, deletedClass.Master.DeletedAt)
	assert.Equal(t, updatedSchoolClass.Master.DeletedBy, deletedClass.Master.DeletedBy)

	// Verify expectations
	mockSchoolClassRepo.AssertExpectations(t)
}

// TestDeleteSchoolClass_NotFound tests the case when the school class is not found
func TestDeleteSchoolClass_NotFound(t *testing.T) {
	// Create mock repositories
	mockSchoolClassRepo := new(MockSchoolClassRepository)
	mockUserSchoolClassRepo := new(MockUserClassRepository)

	// Initialize the service with mock repositories
	schoolClassService := NewSchoolClassService(mockSchoolClassRepo, mockUserSchoolClassRepo)

	// Setup test data
	userID := 1
	classID := uint(1)

	// Define mock behavior
	mockSchoolClassRepo.On("GetSchoolClassByID", classID).Return(models.SchoolClass{}, errors.New("Data not found.")) // Adjusted the error message

	// Call the method under test
	deletedClass, err := schoolClassService.DeleteSchoolClass(classID, userID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, deletedClass)
	assert.Equal(t, "Data not found.", err.Error()) // Adjusted expected error message

	// Verify expectations
	mockSchoolClassRepo.AssertExpectations(t)
}

func TestUpdateSchoolClass_Success(t *testing.T) {
	mockSchoolClassRepo := new(MockSchoolClassRepository)
	mockUserSchoolClassRepo := new(MockUserClassRepository)

	schoolClassService := NewSchoolClassService(mockSchoolClassRepo, mockUserSchoolClassRepo)

	userID := 1
	classID := uint(1)

	mockRequest := &request.SchoolClassCreateUpdateRequest{
		SchoolGradeID:   1,
		PrefixClassID:   1,
		SchoolMajorID:   1,
		Suffix:          "tes",
		SchoolClassName: "tes1",
	}

	mockSchoolClass := models.SchoolClass{
		Master:          models.Master{ID: classID},
		SchoolGradeID:   uint(mockRequest.SchoolGradeID),
		PrefixClassID:   mockRequest.PrefixClassID,
		SchoolMajorID:   mockRequest.SchoolMajorID,
		Suffix:          mockRequest.Suffix,
		SchoolClassName: mockRequest.SchoolClassName,
	}

	// Set expectations untuk GetSchoolClassByID
	mockSchoolClassRepo.On("GetSchoolClassByID", classID).Return(mockSchoolClass, nil)

	mockSchoolClassRepo.On("UpdateSchoolClass", mock.AnythingOfType("*models.SchoolClass")).Return(&mockSchoolClass, nil)

	updatedClass, err := schoolClassService.UpdateSchoolClass(classID, mockRequest, userID)

	assert.NoError(t, err)
	assert.NotNil(t, updatedClass)
	assert.Equal(t, mockRequest.SchoolGradeID, int(updatedClass.SchoolGradeID))
	assert.Equal(t, mockRequest.PrefixClassID, updatedClass.PrefixClassID)
	assert.Equal(t, mockRequest.SchoolMajorID, updatedClass.SchoolMajorID)
	assert.Equal(t, mockRequest.Suffix, updatedClass.Suffix)
	assert.Equal(t, mockRequest.SchoolClassName, updatedClass.SchoolClassName)

	mockSchoolClassRepo.AssertExpectations(t)
}

func TestDetailSchoolClass(t *testing.T) {
	mockRepo := new(MockSchoolClassRepository)
	mockUserRepo := new(MockUserRepository)

	schoolClassService := NewSchoolClassService(mockRepo, mockUserRepo)

	t.Run("Success", func(t *testing.T) {
		mockData := models.SchoolClass{
			Master:          models.Master{ID: 12},
			SchoolGradeID:   3,
			PrefixClassID:   1,
			SchoolMajorID:   1,
			Suffix:          "1",
			SchoolClassName: "XI IPA 4",
		}

		mockRepo.On("GetSchoolClassByID", mock.Anything).Return(mockData, nil)

		// Tangkap kedua nilai yang dikembalikan
		result, err := schoolClassService.GetSchoolClassByID(12)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "XI IPA 4", result.SchoolClassName)
		mockRepo.AssertExpectations(t)
	})
}
