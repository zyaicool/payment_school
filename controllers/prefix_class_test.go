package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPrefixClassService struct {
	mock.Mock
}

func (m *MockPrefixClassService) GetAllPrefixClassService(search string, userID int) (response.PrefixCLassResponse, error) {
	args := m.Called(search, userID)
	return args.Get(0).(response.PrefixCLassResponse), args.Error(1)
}

func (m *MockPrefixClassService) CreatePrefixClassService(prefixClass *request.PrefixClassCreate, userID int) (*models.PrefixClass, error) {
	args := m.Called(prefixClass, userID)
	return args.Get(0).(*models.PrefixClass), args.Error(1)
}

func TestGetPrefixClass(t *testing.T) {
	app := fiber.New()
	mockService := new(MockPrefixClassService)
	controller := &PrefixClassController{
		prefixClassService: mockService,
	}

	tests := []struct {
		name         string
		search       string
		userID       float64
		mockResponse response.PrefixCLassResponse
		mockError    error
		expectedCode int
	}{
		{
			name:   "Success",
			search: "test",
			userID: 1,
			mockResponse: response.PrefixCLassResponse{
				Data: []response.DetailPrefixResponse{
					{ID: 1, PrefixName: "Test Prefix"},
				},
			},
			mockError:    nil,
			expectedCode: fiber.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock JWT middleware
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("user", jwt.MapClaims{"user_id": tc.userID})
				return c.Next()
			})

			mockService.On("GetAllPrefixClassService", tc.search, int(tc.userID)).Return(tc.mockResponse, tc.mockError)

			app.Get("/test", controller.GetPrefixClass)
			req := httptest.NewRequest("GET", "/test?search="+tc.search, nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreatePrefixClass(t *testing.T) {
	app := fiber.New()
	mockService := new(MockPrefixClassService)
	controller := &PrefixClassController{
		prefixClassService: mockService,
	}

	tests := []struct {
		name         string
		input        request.PrefixClassCreate
		userID       float64
		mockResponse *models.PrefixClass
		mockError    error
		expectedCode int
	}{
		{
			name: "Success",
			input: request.PrefixClassCreate{
				PrefixName: "Test",
				SchoolID:   1,
			},
			userID: 1,
			mockResponse: &models.PrefixClass{
				PrefixName: "Test",
				SchoolID:   1,
			},
			mockError:    nil,
			expectedCode: fiber.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app.Use(func(c *fiber.Ctx) error {
				c.Locals("user", jwt.MapClaims{"user_id": tc.userID})
				return c.Next()
			})

			mockService.On("CreatePrefixClassService", &tc.input, int(tc.userID)).Return(tc.mockResponse, tc.mockError)

			app.Post("/test", controller.CreatePrefixClass)

			jsonBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest("POST", "/test", strings.NewReader(string(jsonBody)))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedCode, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}
