package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProvinceService is a mock of ProvinceService
type MockProvinceService struct {
	mock.Mock
}

func (m *MockProvinceService) GetDataProvinces(filename string) ([]services.Region, error) {
	args := m.Called(filename)
	return args.Get(0).([]services.Region), args.Error(1)
}

func TestGetDataProvinces_Success(t *testing.T) {
	// Initialize Fiber app and mock service
	app := fiber.New()
	mockService := new(MockProvinceService)
	controller := NewProvinceController(mockService)

	// Define mock response
	mockRegions := []services.Region{
		{ID: "1", Name: "Province A"},
		{ID: "2", Name: "Province B"},
	}
	mockService.On("GetDataProvinces", "data/mst_province.json").Return(mockRegions, nil)

	// Define route and perform test
	app.Get("/provinces", controller.GetDataProvinces)
	// req, _ := app.Test(fiber.NewRequest("GET", "/provinces"))

	req := httptest.NewRequest(http.MethodGet, "/provinces", nil)

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string][]services.Region
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, mockRegions, response["data"])
}

// func TestGetDataProvinces_Failure(t *testing.T) {
// 	// Initialize Fiber app and mock service
// 	app := fiber.New()
// 	mockService := new(MockProvinceService)
// 	controller := NewProvinceController(mockService)

// 	// Set mock to return an error
// 	mockService.On("GetDataProvinces", "data/mst_province.json").Return(([]services.Region)(nil), errors.New("file not found"))

// 	// Define route and perform test
// 	app.Get("/provinces", controller.GetDataProvinces)
// 	req := httptest.NewRequest(http.MethodGet, "/provinces", nil)
// 	resp, err := app.Test(req, -1)

// 	assert.NoError(t, err)
// 	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

// 	var response map[string]string
// 	json.NewDecoder(resp.Body).Decode(&response)
// 	assert.Equal(t, "Failed to retrieve provinces", response["error"])
// 	assert.Equal(t, "file not found", response["details"])
// }
