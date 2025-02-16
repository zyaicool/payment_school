package controllers

import (
	"fmt"
	"net/http/httptest"
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolYearService struct {
	mock.Mock
}

func (m *MockSchoolYearService) GetSchoolYearByID(id uint) (*response.SchoolYearDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*response.SchoolYearDetailResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSchoolYearService) CreateSchoolYear(req *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error) {
	args := m.Called(req, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SchoolYear), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSchoolYearService) DeleteSchoolYear(id uint, userID int) (*models.SchoolYear, error) {
	args := m.Called(id, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SchoolYear), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSchoolYearService) GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, userID int) (response.SchoolYearListResponse, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder, userID)
	return args.Get(0).(response.SchoolYearListResponse), args.Error(1)
}

func (m *MockSchoolYearService) UpdateSchoolYear(id uint, schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error) {
	args := m.Called(schoolYearRequest, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SchoolYear), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetRoleByID(roleID uint) (*models.Role, error) {
	args := m.Called(roleID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func parseTime(timeString string) *time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeString)
	return &parsedTime
}

func TestGetDetailSchoolYear(t *testing.T) {
	app := fiber.New()
	mockService := new(MockSchoolYearService)
	controller := &SchoolYearController{
		schoolYearService: mockService,
	}

	tests := []struct {
		name         string
		id           int
		mockResponse *response.SchoolYearDetailResponse
		mockError    error
		expectedCode int
	}{
		{
			name: "Success",
			id:   1,
			mockResponse: &response.SchoolYearDetailResponse{
				ID:             1,
				SchoolYearName: "2026/2027",
				StartDate:      parseTime("2026-11-04T07:00:00Z"),
				EndDate:        parseTime("2027-11-29T07:00:00Z"),
			},
			mockError:    nil,
			expectedCode: fiber.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("GetSchoolYearByID", uint(tc.id)).Return(tc.mockResponse, tc.mockError)

			// Update the route to include the ID parameter
			app.Get("/test/:id", controller.GetDataSchoolYear)
			req := httptest.NewRequest("GET", fmt.Sprintf("/test/%d", tc.id), nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateSchoolYear(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		requestBody    string
		userID         int
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Success create school year",
			url:            "/api/v1/schoolYear/create",
			requestBody:    `{"schoolYearName": "2024/2025", "schoolId": 1, "startDate": "2024-01-02", "endDate": "2025-01-08"}`,
			userID:         1,
			expectedStatus: 200,
			mockSetup:      func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest("POST", tc.url, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
		})
	}
}
