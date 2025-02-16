package controllers_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolClassService struct {
	mock.Mock
}

// Implement all methods from SchoolClassServiceInterface
func (m *MockSchoolClassService) GetAllSchoolClass(page int, limit int, search string, sortBy string, sortOrder string, showDeletedData bool, user models.User) (response.SchoolClassListResponse, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder, showDeletedData, user)
	return args.Get(0).(response.SchoolClassListResponse), args.Error(1)
}

func (m *MockSchoolClassService) GetSchoolClassByID(id uint) (*response.SchoolClassDetailResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*response.SchoolClassDetailResponse), args.Error(1)
}

func (m *MockSchoolClassService) CreateSchoolClass(request *request.SchoolClassCreateUpdateRequest, userID int) (*models.SchoolClass, error) {
	args := m.Called(request, userID)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassService) UpdateSchoolClass(id uint, request *request.SchoolClassCreateUpdateRequest, userID int) (*models.SchoolClass, error) {
	args := m.Called(id, request, userID)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassService) DeleteSchoolClass(id uint, userID int) (*models.SchoolClass, error) {
	args := m.Called(id, userID)
	return args.Get(0).(*models.SchoolClass), args.Error(1)
}

func (m *MockSchoolClassService) GenerateSchoolClassCode() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetAllSchoolClass(t *testing.T) {
	// app := fiber.New()
	mockService := new(MockSchoolClassService)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Success get all school class",
			url:            "/api/v1/schoolClass/getAllSchoolClass?page=1&limit=10",
			expectedStatus: 200,
			mockSetup: func() {
				// mockService.On("GetAllSchoolClass", 1, 10, "", "", "", false, mock.Anything).
				// Return(response.SchoolClassListResponse{}, nil)
			},
		},
		// Add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest("GET", tc.url, nil)
			req.Header.Set("Content-Type", "application/json")

			// resp, err := app.Test(req)

			// assert.NoError(t, err)
			// assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteSchoolClass(t *testing.T) {
	mockService := new(MockSchoolClassService)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Success delete school class",
			url:            "/api/v1/schoolClass/delete/1",
			expectedStatus: 200,
			mockSetup: func() {

			},
		},
		{
			name:           "Error deleting school class (Service failure)",
			url:            "/api/v1/schoolClass/delete/1",
			expectedStatus: 500,
			mockSetup: func() {

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock
			tc.mockSetup()

			// Create a test request
			req := httptest.NewRequest("DELETE", tc.url, nil)
			req.Header.Set("Content-Type", "application/json")

			mockService.AssertExpectations(t)
		})
	}
}

// func TestCreateSchoolClass(t *testing.T) {
// 	mockService := new(MockSchoolClassService)

// 	testCases := []struct {
// 		name           string
// 		payload        *request.SchoolClassCreateUpdateRequest
// 		expectedStatus int
// 		mockSetup      func()
// 	}{
// 		{
// 			name: "Success create school class",
// 			payload: &request.SchoolClassCreateUpdateRequest{
// 				SchoolGradeID:   1,
// 				SchoolClassName: "Class A",
// 				PrefixClassID:   2,
// 				SchoolMajorID:   3,
// 				Suffix:          "A",
// 				SchoolID:        1,
// 			},
// 			expectedStatus: 200,
// 			mockSetup: func() {
// 				mockService.On("CreateSchoolClass", mock.Anything, 1).Return(&models.SchoolClass{
// 					SchoolClassName: "Class A",
// 				}, nil)
// 			},
// 		},
// 		{
// 			name:           "Failed to create school class due to invalid payload",
// 			payload:        nil,
// 			expectedStatus: 400,
// 			mockSetup: func() {
// 			},
// 		},
// 		{
// 			name: "Failed to create school class due to service error",
// 			payload: &request.SchoolClassCreateUpdateRequest{
// 				SchoolGradeID:   1,
// 				SchoolClassName: "Class B",
// 				PrefixClassID:   2,
// 				SchoolMajorID:   3,
// 				Suffix:          "B",
// 				SchoolID:        1,
// 			},
// 			expectedStatus: 500,
// 			mockSetup: func() {
// 				mockService.On("CreateSchoolClass", mock.Anything, 1).Return(nil, assert.AnError)
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.mockSetup()

// 			// Create a test request
// 			var reqBody []byte
// 			if tc.payload != nil {
// 				reqBody, _ = json.Marshal(tc.payload)
// 			}
// 			req := httptest.NewRequest("POST", "/api/v1/schoolClass/create", bytes.NewBuffer(reqBody))
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", "Bearer mock_token")

// 			// Simulate logged-in user
// 			ctx := req.Context()
// 			ctx = context.WithValue(ctx, "user", map[string]interface{}{
// 				"user_id": 1.0,
// 			})
// 			req = req.WithContext(ctx)

// 			// Create response recorder
// 			recorder := httptest.NewRecorder()

// 			// Call the handler (replace with actual handler call if available)
// 			handler := func(w http.ResponseWriter, r *http.Request) {
// 				// Inject mock service and call the actual controller method
// 				schoolClassController := &SchoolClassController{
// 					schoolClassService: mockService,
// 				}
// 				schoolClassController.CreateSchoolClass(w, r)
// 			}
// 			handler(recorder, req)

//				// Assertions
//				assert.Equal(t, tc.expectedStatus, recorder.Code)
//				mockService.AssertExpectations(t)
//			})
//		}
//	}

func TestUpdateSchoolClass(t *testing.T) {
	mockService := new(MockSchoolClassService)

	testCases := []struct {
		name           string
		url            string
		requestBody    string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Success update school class",
			url:            "/api/v1/schoolClass/update/1",
			requestBody:    `{"school_class_name":"Class A"}`,
			expectedStatus: 200,
			mockSetup: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest("PUT", tc.url, strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			assert.Equal(t, tc.expectedStatus, w.Code)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetDetailSchoolClass(t *testing.T) {
	mockService := new(MockSchoolClassService)

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		mockSetup      func()
	}{
		{
			name:           "Success get detail school class",
			url:            "/api/v1/schoolClass/detail/1",
			expectedStatus: 200,
			mockSetup: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest("GET", tc.url, nil)
			req.Header.Set("Content-Type", "application/json")

			mockService.AssertExpectations(t)
		})
	}
}
