package services

import (
	"errors"
	request "schoolPayment/dtos/request"
	"schoolPayment/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolYearRepository struct {
	mock.Mock
}

func (m *MockSchoolYearRepository) GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, schoolId uint) ([]models.SchoolYearList, int64, int, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder, schoolId)
	return args.Get(0).([]models.SchoolYearList),
		args.Get(1).(int64),
		args.Get(2).(int),
		args.Error(3)
}

func (m *MockSchoolYearRepository) CreateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error) {
	args := m.Called(schoolYear)
	return args.Get(0).(*models.SchoolYear), args.Error(1)
}

func (m *MockSchoolYearRepository) UpdateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error) {
	args := m.Called(schoolYear)
	return args.Get(0).(*models.SchoolYear), args.Error(1)
}

func (m *MockSchoolYearRepository) GetSchoolYearByID(id uint) (models.SchoolYear, error) {
	args := m.Called(id)
	if schoolYear, ok := args.Get(0).(*models.SchoolYear); ok {
		return *schoolYear, args.Error(1)
	}
	return models.SchoolYear{}, args.Error(1)
}

func (m *MockSchoolYearRepository) GetAllSchoolYearsBulk(yearNames []string, schoolID uint) ([]models.SchoolYear, error) {
	args := m.Called(yearNames, schoolID)
	return args.Get(0).([]models.SchoolYear), args.Error(1)
}

func (m *MockSchoolYearRepository) GetLastSequenceNumberSchoolYears() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func parseTime(timeString string) *time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeString)
	return &parsedTime
}

func TestGetSchoolYearByIDService(t *testing.T) {
	mockRepo := new(MockSchoolYearRepository)
	mockUserRepo := new(MockUserRepository)

	service := NewSchoolYearService(mockRepo, mockUserRepo)

	testCases := []struct {
		name          string
		id            uint
		mockResponse  models.SchoolYear
		mockError     error
		expectedError error
	}{
		{
			name: "Success",
			id:   28,
			mockResponse: models.SchoolYear{
				SchoolYearName: "2026/2027",
				StartDate:      parseTime("2026-11-04T07:00:00Z"),
				EndDate:        parseTime("2027-11-29T07:00:00Z"),
			},
			mockError:     nil,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetSchoolYearByID", tc.id).Return(&tc.mockResponse, tc.mockError).Once()
			result, err := service.GetSchoolYearByID(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockResponse.SchoolYearName, result.SchoolYearName)

				if tc.mockResponse.StartDate != nil && result.StartDate != nil {
					assert.Equal(t, *tc.mockResponse.StartDate, *result.StartDate)
				}
				if tc.mockResponse.EndDate != nil && result.EndDate != nil {
					assert.Equal(t, *tc.mockResponse.EndDate, *result.EndDate)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateSchoolYearService(t *testing.T) {
	mockSchoolYearRepo := new(MockSchoolYearRepository)
	mockUserRepo := new(MockUserRepository)

	// Debug: pastikan service dibuat dengan benar
	t.Logf("Creating service with repos: %v, %v", mockSchoolYearRepo, mockUserRepo)
	service := NewSchoolYearService(mockSchoolYearRepo, mockUserRepo)

	testCases := []struct {
		name               string
		request            request.SchoolYearCreateUpdateRequest
		lastSequenceNumber int
		repoError          error
		expectedResult     *models.SchoolYear
		expectedError      error
	}{
		{
			name: "Success",
			request: request.SchoolYearCreateUpdateRequest{
				SchoolYearName: "2023/2025",
				SchoolId:       1,
				StartDate:      "2023-01-01",
				EndDate:        "2025-12-31",
			},
			lastSequenceNumber: 1,
			repoError:          nil,
			expectedResult: &models.SchoolYear{
				Master: models.Master{
					CreatedBy: 1,
					UpdatedBy: 1,
				},
				SchoolYearCode: "SY001",
				SchoolYearName: "2023/2025",
				SchoolId:       1,
				StartDate: func() *time.Time {
					t, _ := time.Parse("2006-01-02", "2023-01-01")
					return &t
				}(),
				EndDate: func() *time.Time {
					t, _ := time.Parse("2006-01-02", "2025-12-31")
					return &t
				}(),
			},
			expectedError: nil,
		},
		{
			name: "Error fetching sequence",
			request: request.SchoolYearCreateUpdateRequest{
				SchoolYearName: "2023/2025",
				SchoolId:       1,
				StartDate:      "2023-01-01",
				EndDate:        "2025-12-31",
			},
			lastSequenceNumber: 0,
			repoError:          errors.New("DB error"),
			expectedResult:     nil,
			expectedError:      errors.New("DB error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSchoolYearRepo.On("GetLastSequenceNumberSchoolYears").Return(tc.lastSequenceNumber, tc.repoError).Once()

			if tc.repoError == nil {
				mockSchoolYearRepo.On("CreateSchoolYear", mock.Anything).Return(tc.expectedResult, nil).Once()
			}

			result, err := service.CreateSchoolYear(&tc.request, 1)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			mockSchoolYearRepo.AssertExpectations(t)
		})
	}
}
