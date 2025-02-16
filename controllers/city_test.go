package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for CityService
type MockCityService struct {
	mock.Mock
}

func (m *MockCityService) GetDataCities(filename string, provinceID int) ([]services.DistrictNewMap, error) {
	args := m.Called(filename, provinceID)
	if args.Get(0) != nil {
		return args.Get(0).([]services.DistrictNewMap), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetDataCities_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockCityService)
	cityController := NewCityController(mockService)

	// Mock data
	mockDistricts := []services.DistrictNewMap{
		{ID: 1, IDProvince: 1, District: "District1"},
		{ID: 2, IDProvince: 1, District: "District2"},
	}
	mockService.On("GetDataCities", "data/mst_district.json", 1).Return(mockDistricts, nil)

	// Define Fiber route
	app.Get("/cities", cityController.GetDataCities)

	req := httptest.NewRequest("GET", "/cities?provinceId=1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.NotNil(t, body["data"])
	assert.Equal(t, len(body["data"].([]interface{})), len(mockDistricts))
	mockService.AssertExpectations(t)
}

func TestGetDataCities_NoDistrictsFound(t *testing.T) {
	app := fiber.New()
	mockService := new(MockCityService)
	cityController := NewCityController(mockService)

	mockService.On("GetDataCities", "data/mst_district.json", 1).Return(nil, errors.New("No districts found for this province"))

	// Define Fiber route
	app.Get("/cities", cityController.GetDataCities)

	req := httptest.NewRequest("GET", "/cities?provinceId=1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "No districts found for this province", body["error"])
	mockService.AssertExpectations(t)
}

func TestGetDataCities_InvalidProvinceID(t *testing.T) {
	app := fiber.New()
	mockService := new(MockCityService)
	cityController := NewCityController(mockService)

	mockService.On("GetDataCities", "data/mst_district.json", 0).Return(nil, errors.New("Invalid province ID"))

	// Define Fiber route
	app.Get("/cities", cityController.GetDataCities)

	req := httptest.NewRequest("GET", "/cities?provinceId=0", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "Invalid province ID", body["error"])
	mockService.AssertExpectations(t)
}
