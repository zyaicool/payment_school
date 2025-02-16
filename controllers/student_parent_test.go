package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"schoolPayment/controllers"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudentParentService struct {
	mock.Mock
}

func (m *MockStudentParentService) GetAllParent(page int, limit int, search string, userID int) (response.StudentParentResponse, error) {
	args := m.Called(page, limit, search, userID)
	return args.Get(0).(response.StudentParentResponse), args.Error(1)
}

func (m *MockStudentParentService) GetParentByID(id uint, userID int) (*models.StudentParent, error) {
	args := m.Called(id, userID)
	return args.Get(0).(*models.StudentParent), args.Error(1)
}

func (m *MockStudentParentService) CreateParent(parent *models.StudentParent, userID int) (*models.StudentParent, error) {
	args := m.Called(parent, userID)
	return args.Get(0).(*models.StudentParent), args.Error(1)
}

func (m *MockStudentParentService) UpdateParent(id uint, parent *models.StudentParent, userID int) (*models.StudentParent, error) {
	args := m.Called(id, parent, userID)
	return args.Get(0).(*models.StudentParent), args.Error(1)
}

func (m *MockStudentParentService) GetParentByUserLogin(id uint) (*models.StudentParent, error) {
	args := m.Called(id)
	return args.Get(0).(*models.StudentParent), args.Error(1)
}

func (m *MockStudentParentService) CreateBatchParent(parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error) {
	args := m.Called(parents, userID)
	return args.Get(0).(*[]*models.StudentParent), args.Error(1)
}

func (m *MockStudentParentService) UpdateBatchParent(id uint, parents []*models.StudentParent, userID int) (*[]*models.StudentParent, error) {
	args := m.Called(id, parents, userID)
	return args.Get(0).(*[]*models.StudentParent), args.Error(1)
}

func TestGetAllParent(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockData := response.StudentParentResponse{
		Data: []models.StudentParent{
			{
				Master: models.Master{ID: 1},
				ParentName: "Parent 1",
			},
			{
				Master: models.Master{ID: 2},
				ParentName: "Parent 2",
			},
		},
		Page:  1,
		Limit: 10,
	}
	mockService.On("GetAllParent", 1, 10, "", 1).Return(mockData, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/studentParent/getAllParent", controller.GetAllParent)

	req := httptest.NewRequest("GET", "/api/v1/studentParent/getAllParent?page=1&limit=10", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetParentByID(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParent := &models.StudentParent{
		Master: models.Master{ID: 1},
		ParentName: "Parent 1",
	}
	mockService.On("GetParentByID", uint(1), 1).Return(mockParent, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/studentParent/detail/:id", controller.GetDataParent)

	req := httptest.NewRequest("GET", "/api/v1/studentParent/detail/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCreateParent(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParent := &models.StudentParent{
		Master: models.Master{ID: 1},
		ParentName: "Parent 1",
	}
	mockService.On("CreateParent", mockParent, 1).Return(mockParent, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/studentParent/create", controller.CreateDataParent)

	body, _ := json.Marshal(mockParent)
	req := httptest.NewRequest("POST", "/api/v1/studentParent/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdateParent(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParent := &models.StudentParent{
		Master: models.Master{ID: 1},
		ParentName: "Updated Parent",
	}
	mockService.On("UpdateParent", uint(1), mockParent, 1).Return(mockParent, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Put("/api/v1/studentParent/update/:id", controller.UpdateParent)

	body, _ := json.Marshal(mockParent)
	req := httptest.NewRequest("PUT", "/api/v1/studentParent/update/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetParentByUserLogin(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParent := &models.StudentParent{
		Master: models.Master{ID: 1},
		ParentName: "Parent 1",
	}
	mockService.On("GetParentByUserLogin", uint(1)).Return(mockParent, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/studentParent/getParentByUserLogin", controller.GetDataParentByUserLogin)

	req := httptest.NewRequest("GET", "/api/v1/studentParent/getParentByUserLogin", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCreateBatchParent(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParents := []*models.StudentParent{
		{
			Master: models.Master{ID: 1},
			ParentName: "Parent 1",
		},
		{
			Master: models.Master{ID: 2},
			ParentName: "Parent 2",
		},
	}
	mockService.On("CreateBatchParent", mockParents, 1).Return(&mockParents, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, 
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/studentParent/createBatch", controller.CreateBatchDataParent)

	body, _ := json.Marshal(mockParents)
	req := httptest.NewRequest("POST", "/api/v1/studentParent/createBatch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdateBatchParent(t *testing.T) {
	mockService := new(MockStudentParentService)

	// Mock Data
	mockParents := []*models.StudentParent{
		{
			Master: models.Master{ID: 1},
			ParentName: "Updated Parent 1",
		},
		{
			Master: models.Master{ID: 2},
			ParentName: "Updated Parent 2",
		},
	}
	mockService.On("UpdateBatchParent", uint(1), mockParents, 1).Return(&mockParents, nil)

	controller := controllers.NewStudentParentController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Put("/api/v1/studentParent/updateBatch/:id", controller.UpdateBatchParent)

	body, _ := json.Marshal(mockParents)
	req := httptest.NewRequest("PUT", "/api/v1/studentParent/updateBatch/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func generateJWTStudentParent(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}