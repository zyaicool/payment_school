package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"schoolPayment/controllers"
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchoolMajorService struct {
	mock.Mock
}

func (m *MockSchoolMajorService) GetAllSchoolMajorService(search string, roleID int) (response.SchoolMajorResponse, error) {
	args := m.Called(search, roleID)
	return args.Get(0).(response.SchoolMajorResponse), args.Error(1)
}

func (m *MockSchoolMajorService) CreateSchoolMajorService(major *request.SchoolMajorCreate, userID int) (*models.SchoolMajor, error) {
	args := m.Called(major, userID)
	return args.Get(0).(*models.SchoolMajor), args.Error(1)
}

func TestGetAllSchoolMajor(t *testing.T) {
	mockService := new(MockSchoolMajorService)

	// Mock Data
	mockData := response.SchoolMajorResponse{
		Data: []response.DetailSchoolMajorResponse{
			{ID: 1, SchoolMajorName: "Physics"},
			{ID: 2, SchoolMajorName: "Mathematics"},
		},
	}
	mockService.On("GetAllSchoolMajorService", "", 1).Return(mockData, nil)

	controller := controllers.NewSchoolMajorController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"role_id": 1.0, 

	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/schoolMajor/getListSchoolMajor", controller.GetAllSchoolMajor)

	req := httptest.NewRequest("GET", "/api/v1/schoolMajor/getListSchoolMajor", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func generateJWT(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

func TestCreateSchoolMajor(t *testing.T) {
	mockService := new(MockSchoolMajorService)
	mockMajor := request.SchoolMajorCreate{
		SchoolMajorName: "test school major",
		SchoolID:        1,
	}
	mockResponse := &models.SchoolMajor{
		SchoolID:        1,
		SchoolMajorName: "test school major",
	}

	mockService.On("CreateSchoolMajorService", &mockMajor, 0).Return(mockResponse, nil)

	controller := controllers.NewSchoolMajorController(mockService)
	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 0, 
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/schoolMajor/create", controller.CreateSchoolMajor)

	body, _ := json.Marshal(mockMajor)
	req := httptest.NewRequest("POST", "/api/v1/schoolMajor/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}
