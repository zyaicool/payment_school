// Package services_test contains unit tests for the services package
package services_test

import (
	"schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPrefixClassRepository is a mock implementation of the PrefixClassRepository interface
// used for testing the PrefixClassService
type MockPrefixClassRepository struct {
	mock.Mock
}

// GetAllPrefixClassRepository mocks the repository method for getting all prefix classes
func (m *MockPrefixClassRepository) GetAllPrefixClassRepository(search string, user models.User) ([]models.PrefixClass, error) {
	args := m.Called(search, user)
	return args.Get(0).([]models.PrefixClass), args.Error(1)
}

// CreatePrefixClassRepository mocks the repository method for creating a prefix class
func (m *MockPrefixClassRepository) CreatePrefixClassRepository(prefix *models.PrefixClass) (*models.PrefixClass, error) {
	args := m.Called(prefix)
	return args.Get(0).(*models.PrefixClass), args.Error(1)
}

// CheckPrefixClassExists mocks the repository method for checking if a prefix class exists
func (m *MockPrefixClassRepository) CheckPrefixClassExists(prefixName string, schoolID uint) (bool, error) {
	args := m.Called(prefixName, schoolID)
	return args.Bool(0), args.Error(1)
}

// MockUserRepository is a mock implementation of the UserRepository interface
// used for testing user-related functionality in the PrefixClassService
type MockUserRepository struct {
	mock.Mock
}

// GetUserByID mocks the repository method for getting a user by ID
func (m *MockUserRepository) GetUserByID(id uint) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
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

// TestGetAllPrefixClassService tests the GetAllPrefixClassService method
// It verifies that the service correctly retrieves all prefix classes based on search criteria
// and user permissions
func TestGetAllPrefixClassService(t *testing.T) {
	mockPrefixRepo := new(MockPrefixClassRepository)
	mockUserRepo := new(MockUserRepository)

	service := services.NewPrefixClassServiceWithRepo(mockPrefixRepo, mockUserRepo)

	testCases := []struct {
		name          string
		search        string
		userID        int
		mockUser      models.User
		mockResponse  []models.PrefixClass
		mockError     error
		expectedError error
	}{
		{
			name:   "Success get all prefix classes",
			search: "",
			userID: 1,
			mockUser: models.User{
				UserSchool: &models.UserSchool{
					SchoolID: 1,
				},
			},
			mockResponse: []models.PrefixClass{
				{Master: models.Master{}, PrefixName: "A", SchoolID: 1},
				{Master: models.Master{}, PrefixName: "B", SchoolID: 1},
			},
			mockError:     nil,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.On("GetUserByID", uint(tc.userID)).Return(tc.mockUser, nil).Once()
			mockPrefixRepo.On("GetAllPrefixClassRepository", tc.search, tc.mockUser).
				Return(tc.mockResponse, tc.mockError).Once()

			result, err := service.GetAllPrefixClassService(tc.search, tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.mockResponse), len(result.Data))
				for i, prefix := range result.Data {
					assert.Equal(t, tc.mockResponse[i].PrefixName, prefix.PrefixName)
				}
			}

			mockUserRepo.AssertExpectations(t)
			mockPrefixRepo.AssertExpectations(t)
		})
	}
}

// TestCreatePrefixClassService tests the CreatePrefixClassService method
// It verifies that the service correctly handles the creation of new prefix classes,
// including validation of duplicate prefix names and proper error handling
func TestCreatePrefixClassService(t *testing.T) {
	mockRepo := new(MockPrefixClassRepository)
	mockUserRepo := new(MockUserRepository)

	service := services.NewPrefixClassServiceWithRepo(mockRepo, mockUserRepo)

	testCases := []struct {
		name          string
		request       request.PrefixClassCreate
		userID        int
		user          models.User
		existsResult  bool
		existsError   error
		createResult  *models.PrefixClass
		createError   error
		expectedError error
	}{
		{
			name: "Success create prefix class",
			request: request.PrefixClassCreate{
				PrefixName: "X",
				SchoolID:   1,
			},
			userID: 1,
			user: models.User{
				UserSchool: &models.UserSchool{
					SchoolID: 1,
				},
			},
			existsResult: false,
			existsError:  nil,
			createResult: &models.PrefixClass{
				Master:     models.Master{},
				PrefixName: "X",
				SchoolID:   1,
			},
			createError:   nil,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("CheckPrefixClassExists", tc.request.PrefixName, tc.request.SchoolID).
				Return(tc.existsResult, tc.existsError).Once()

			if tc.existsError == nil && !tc.existsResult {
				expectedPrefix := &models.PrefixClass{
					Master: models.Master{
						CreatedBy: tc.userID,
					},
					PrefixName: tc.request.PrefixName,
					SchoolID:   tc.request.SchoolID,
				}
				mockRepo.On("CreatePrefixClassRepository", expectedPrefix).
					Return(tc.createResult, tc.createError).Once()
			}

			result, err := service.CreatePrefixClassService(&tc.request, tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.createResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
