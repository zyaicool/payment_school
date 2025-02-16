package services

import (
	"errors"
	"fmt"
	"testing"
	"time"

	request "schoolPayment/dtos/request"
	"schoolPayment/models"

	"github.com/stretchr/testify/assert"
)

// MockBillingTypeRepository is a manual mock for the BillingTypeRepository
type MockBillingTypeRepository struct {
	// You can define fields to store expected data if needed
	billingTypeCode string
	billingTypeList []models.BillingType
	billingType     models.BillingType
	userSchool      models.UserSchool
	lastSequence    int
	err             error
}

// GetAllBillingType returns the mock billing type list and error
func (m *MockBillingTypeRepository) GetAllBillingType(page int, limit int, search string) ([]models.BillingType, error) {
	return m.billingTypeList, m.err
}

func (m *MockBillingTypeRepository) GetBillingTypeByIDTest(id uint) (models.BillingType, error) {
	if m.err != nil {
		return models.BillingType{}, m.err
	}
	return m.billingType, nil
}

// Mock UpdateBillingType function
func (m *MockBillingTypeRepository) UpdateBillingType(billingType *models.BillingType) (*models.BillingType, error) {
	if m.err != nil {
		return nil, m.err
	}
	return billingType, nil
}

// GetUserSchoolByUserId simulates fetching user school by user ID
func (m *MockBillingTypeRepository) GetUserSchoolByUserIdTest(userID uint) (models.UserSchool, error) {
	return m.userSchool, m.err
}

// CreateBillingType simulates creating a new billing type
func (m *MockBillingTypeRepository) CreateBillingType(billingType *models.BillingType) (*models.BillingType, error) {
	return &m.billingType, m.err
}

func (m *MockBillingTypeRepository) GetLastSequenceNumberBillingType() (int, error) {
	return m.lastSequence, m.err
}

func TestGetAllBillingType(t *testing.T) {
	// Define the table of test cases
	tests := []struct {
		name            string
		page, limit     int
		search          string
		billingTypeList []models.BillingType
		expectedLength  int
		expectedError   error
	}{
		{
			name: "Returns all billing types",
			page: 1, limit: 10, search: "",
			billingTypeList: []models.BillingType{
				{SchoolID: 1, BillingTypeName: "Basic"},
				{SchoolID: 2, BillingTypeName: "Premium"},
			},
			expectedLength: 2,
			expectedError:  nil,
		},
		{
			name: "Returns empty list when no match",
			page: 1, limit: 10, search: "asdasdasd",
			billingTypeList: []models.BillingType{},
			expectedLength:  0,
			expectedError:   nil,
		},
	}

	// Loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockBillingTypeRepository{
				billingTypeList: tt.billingTypeList,
				err:             tt.expectedError,
			}

			// Call service and perform test
			service := NewBillingTypeService(mockRepo)
			result, err := service.GetAllBillingType(tt.page, tt.limit, tt.search)

			// Validate results
			assert.Equal(t, tt.expectedError, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedLength, len(result.Data))
		})
	}
}

func TestGetBillingTypeByID(t *testing.T) {
	// Define the test cases in a table
	testCases := []struct {
		name           string
		billingType    models.BillingType
		err            error
		inputID        uint
		expectedError  string
		expectedResult models.BillingType
	}{
		{
			name:           "BillingType Found",
			billingType:    models.BillingType{SchoolID: 1, BillingTypeName: "Basic"},
			err:            nil,
			inputID:        1,
			expectedError:  "",
			expectedResult: models.BillingType{SchoolID: 1, BillingTypeName: "Basic"},
		},
		{
			name:           "BillingType Not Found",
			billingType:    models.BillingType{},
			err:            errors.New("record not found"),
			inputID:        999,
			expectedError:  "record not found",
			expectedResult: models.BillingType{},
		},
	}

	// Iterate over the test cases and run them
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository with test case data
			mockRepo := &MockBillingTypeRepository{
				billingType: tc.billingType,
				err:         tc.err,
			}

			// Call the service method
			service := NewBillingTypeService(mockRepo)
			result, err := service.GetBillingTypeByID(tc.inputID)

			// Validate results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestUpdateBillingType(t *testing.T) {
	testCases := []struct {
		name                string
		billingType         *models.BillingType
		billingTypeRequest  *request.BillingTypeCreateUpdateRequest
		userID              int
		err                 error
		expectedError       string
		expectedBillingType *models.BillingType
	}{
		{
			name: "Successful Update",
			billingType: &models.BillingType{
				BillingTypeName:   "Old Name",
				BillingTypePeriod: "asd",
				Master: models.Master{
					UpdatedBy: 0,
				},
			},
			billingTypeRequest: &request.BillingTypeCreateUpdateRequest{
				BillingTypeName:   "New Name",
				BillingTypePeriod: "asd",
			},
			userID:        1,
			err:           nil,
			expectedError: "",
			expectedBillingType: &models.BillingType{
				BillingTypeName:   "New Name",
				BillingTypePeriod: "asd",
				Master: models.Master{
					UpdatedBy: 1,
				},
			},
		},
		{
			name:                "Billing Type Not Found",
			billingType:         &models.BillingType{},
			billingTypeRequest:  &request.BillingTypeCreateUpdateRequest{},
			userID:              1,
			err:                 errors.New("Data not found."),
			expectedError:       "Data not found.",
			expectedBillingType: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock repository with test case data
			mockRepo := &MockBillingTypeRepository{
				billingType: *tc.billingType,
				err:         tc.err,
			}

			// Call the UpdateBillingType function
			// Call the service method
			service := NewBillingTypeService(mockRepo)
			updatedBillingType, err := service.UpdateBillingType(1, tc.billingTypeRequest, tc.userID)

			// Validate the result
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBillingType, updatedBillingType)
			}
		})
	}
}

func TestDeleteBillingType(t *testing.T) {
	hardcodedTime := time.Date(2023, 10, 24, 12, 0, 0, 0, time.UTC)
	hardcodedTimePointer := &hardcodedTime

	testCases := []struct {
		name                string
		billingType         *models.BillingType
		userID              int
		err                 error
		expectedError       string
		expectedBillingType *models.BillingType
	}{
		{
			name: "Successful Deletion",
			billingType: &models.BillingType{
				BillingTypeName: "Basic",
				Master: models.Master{
					DeletedAt: nil,
					DeletedBy: nil,
				},
			},
			userID:        1,
			err:           nil,
			expectedError: "",
			expectedBillingType: &models.BillingType{
				BillingTypeName: "Basic",
				Master: models.Master{
					DeletedAt: hardcodedTimePointer,
					DeletedBy: func() *int {
						userID := 1
						return &userID
					}(),
				},
			},
		},
		{
			name:                "Billing Type Not Found",
			billingType:         &models.BillingType{},
			userID:              1,
			err:                 errors.New("Data not found."),
			expectedError:       "Data not found.",
			expectedBillingType: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock repository with test case data
			mockRepo := &MockBillingTypeRepository{
				billingType: *tc.billingType,
				err:         tc.err,
			}

			// Call the DeleteBillingType function
			// Call the service method
			service := NewBillingTypeService(mockRepo)
			deletedBillingType, err := service.DeleteBillingType(1, tc.userID)

			// Validate the result
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, deletedBillingType)
			}
		})
	}
}

func TestGenerateBillingTypeCode(t *testing.T) {
	testCases := []struct {
		name          string
		lastSequence  int
		expectedCode  string
		mockError     error
		expectedError string
	}{
		{
			name:          "Generate new code with valid sequence",
			lastSequence:  10,
			expectedCode:  "BT011",
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "Generate new code with zero sequence",
			lastSequence:  0,
			expectedCode:  "BT001",
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "Error while retrieving last sequence",
			lastSequence:  0,
			expectedCode:  "",
			mockError:     errors.New("Database error"),
			expectedError: "Database error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository with the last sequence number and error
			mockRepo := &MockBillingTypeRepository{
				lastSequence: tc.lastSequence,
				err:          tc.mockError,
			}

			// Call the service function to generate the code
			service := NewBillingTypeService(mockRepo)
			code, err := service.GenerateBillingTypeCode()

			// Validate the result
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, code)
			}
		})
	}
}

func TestCreateBillingType(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name                string
		billingTypeRequest  *request.BillingTypeCreateUpdateRequest
		userID              int
		mockBillingTypeCode string
		mockUserSchool      models.UserSchool
		mockBillingType     models.BillingType
		mockError           error
		expectedError       string
		expectedResult      *models.BillingType
	}{
		{
			name: "Success - Generate Code and Create BillingType",
			billingTypeRequest: &request.BillingTypeCreateUpdateRequest{
				SchoolID:          0,
				BillingTypeName:   "Standard",
				BillingTypePeriod: "Monthly",
			},
			userID:              1,
			mockBillingTypeCode: "BT001",
			mockUserSchool:      models.UserSchool{SchoolID: 1},
			mockBillingType: models.BillingType{
				SchoolID:          1,
				BillingTypeCode:   "BT001",
				BillingTypeName:   "Standard",
				BillingTypePeriod: "Monthly",
			},
			mockError:      nil,
			expectedError:  "",
			expectedResult: &models.BillingType{SchoolID: 1, BillingTypeCode: "BT001", BillingTypeName: "Standard", BillingTypePeriod: "Monthly"},
		},
		{
			name: "Error - Failed to Generate BillingType Code",
			billingTypeRequest: &request.BillingTypeCreateUpdateRequest{
				SchoolID:          0,
				BillingTypeName:   "Standard",
				BillingTypePeriod: "Monthly",
			},
			userID:              1,
			mockBillingTypeCode: "",
			mockError:           fmt.Errorf("failed to generate billing type code"),
			expectedError:       "failed to generate billing type code",
			expectedResult:      nil,
		},
		{
			name: "Error - Failed to Fetch User School",
			billingTypeRequest: &request.BillingTypeCreateUpdateRequest{
				SchoolID:          0,
				BillingTypeName:   "Standard",
				BillingTypePeriod: "Monthly",
			},
			userID:              1,
			mockBillingTypeCode: "BT001",
			mockUserSchool:      models.UserSchool{},
			mockError:           fmt.Errorf("failed to fetch user school"),
			expectedError:       "failed to fetch user school",
			expectedResult:      nil,
		},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository with expected behaviors
			mockRepo := &MockBillingTypeRepository{
				billingTypeCode: tc.mockBillingTypeCode,
				userSchool:      tc.mockUserSchool,
				billingType:     tc.mockBillingType,
				err:             tc.mockError,
			}

			// Create the service with the mock repository
			service := BillingTypeService{
				billingRepo: mockRepo,
			}

			// Call the CreateBillingType method
			result, err := service.CreateBillingType(tc.billingTypeRequest, tc.userID)

			// Validate results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedResult.BillingTypeCode, result.BillingTypeCode)
				assert.Equal(t, tc.expectedResult.SchoolID, result.SchoolID)
				assert.Equal(t, tc.expectedResult.BillingTypeName, result.BillingTypeName)
				assert.Equal(t, tc.expectedResult.BillingTypePeriod, result.BillingTypePeriod)
			}
		})
	}
}
