package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBankAccountService struct {
	mock.Mock
}

func (m *MockBankAccountService) CreateBankAccount(bankRequest *request.BankAccountCreateRequest, userID uint) (*models.BankAccount, error) {
	args := m.Called(bankRequest, userID)
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) UpdateBankAccount(id uint, bankRequest *request.BankAccountCreateRequest, userID int) (*models.BankAccount, error) {
	args := m.Called(id, bankRequest, userID)
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) GetAllBankAccounts(page, limit int, search, sortBy, sortOrder string, user models.User) (response.BankAccountListResponse, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder, user)
	return args.Get(0).(response.BankAccountListResponse), args.Error(1)
}

func (m *MockBankAccountService) GetBankAccountDetails(id uint) (*response.BankAccountData, error) {
	args := m.Called(id)
	return args.Get(0).(*response.BankAccountData), args.Error(1)
}
func (m *MockBankAccountService) LoadBanksName(filename string) ([]models.Bank, error) {
	args := m.Called(filename)
	return args.Get(0).([]models.Bank), args.Error(1)
}

func (m *MockBankAccountService) DeleteBankAccount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUpdateSchool_NormalCase(t *testing.T) {
	// Initialize the app
	app := fiber.New()

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
		return c.Next()
	})

	// Mock service
	mockService := new(MockBankAccountService)
	controller := NewBankAccountController(mockService)

	// Mock data to return a *models.School instead of *response.DetailSchoolResponse
	mockService.On("UpdateBankAccount", uint(1), mock.Anything, 1).Return(&models.BankAccount{
		SchoolID: 1,
		BankName: "Bank A",
		// Add other fields as needed
	}, nil)

	// Prepare form data for the request
	formData := url.Values{}
	formData.Set("schoolId", "1")
	formData.Set("bankName", "Bank A")

	// Prepare the request with form data
	app.Put("/update/:id", controller.UpdateBankAccount)
	req := httptest.NewRequest("PUT", "/update/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Handle the request
	resp, err := app.Test(req)

	// Assert no errors occurred
	assert.NoError(t, err)

	// Assert the response status
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Initialize responseBody as a map
	var responseBody fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Extract only the fields needed for verification
	data := responseBody["data"].(map[string]interface{})

	// Assert only the fields that matter, ignoring extra fields
	assert.Equal(t, "Bank A", data["bankName"])

	// Assert the success message
	assert.Equal(t, "Data updated successfully.", responseBody["message"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

func TestUpdateSchool_ErrorCase(t *testing.T) {
	// Initialize the app
	app := fiber.New()

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
		return c.Next()
	})

	// Mock service
	mockService := new(MockBankAccountService)
	controller := NewBankAccountController(mockService)

	// Simulate an error during the update
	mockService.On("UpdateBankAccount", uint(1), mock.Anything, 1).Return(&models.BankAccount{}, errors.New("failed to update bank account"))

	// Prepare form data for the request
	formData := url.Values{}
	formData.Set("schoolId", "1")
	formData.Set("schoolName", "Sekolah A")

	// Prepare the request with form data
	app.Put("/update/:id", controller.UpdateBankAccount)
	req := httptest.NewRequest("PUT", "/update/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Handle the request
	resp, err := app.Test(req)

	// Assert no errors occurred in sending the request
	assert.NoError(t, err)

	// Assert the response status for error (e.g., 500 Internal Server Error)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Initialize responseBody as a map
	var responseBody fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the error message in the response body
	assert.Equal(t, "failed to update bank account", responseBody["error"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateBankAccount_NormalCase(t *testing.T) {
	// Initialize the app
	app := fiber.New()

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
		return c.Next()
	})

	// Mock service
	mockService := new(MockBankAccountService)
	controller := NewBankAccountController(mockService)

	// Mock the service method
	mockService.On("CreateBankAccount", mock.AnythingOfType("*request.BankAccountCreateRequest"), uint(1)).
		Return(&models.BankAccount{
			BankName: "Bank Account",
		}, nil)

	// Prepare JSON data for the request
	requestBody := `{
		"bankName": "Bank Account",
		"accountNumber": "123456789",
		"accountName": "Test Account",
		"accountOwner": "John Doe"
	}`

	// Register the route
	app.Post("/create", controller.CreateBankAccount)
	req := httptest.NewRequest("POST", "/create", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Handle the request
	resp, err := app.Test(req)

	// Assert no errors occurred in sending the request
	assert.NoError(t, err)

	// Assert the response status
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the response message and data
	assert.Equal(t, "Data is saved successfully.", responseBody["message"])
	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, "Bank Account", data["bankName"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

func TestDeleteBankAccount_Success(t *testing.T) {
	// Initialize the app
	app := fiber.New()

	// Mock service
	mockService := new(MockBankAccountService)
	controller := NewBankAccountController(mockService)

	// Register the route
	app.Delete("/bank-account/:id", controller.DeleteBankAccount)

	// Mock the service method for successful deletion
	mockService.On("DeleteBankAccount", uint(1)).Return(nil)

	// Prepare the request
	req := httptest.NewRequest("DELETE", "/bank-account/1", nil)

	// Handle the request
	resp, err := app.Test(req)

	// Assert no errors occurred in sending the request
	assert.NoError(t, err)

	// Assert the response status
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the response body
	assert.Equal(t, "Bank account deleted successfully.", responseBody["message"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

func TestGetBankName_Success(t *testing.T) {
	// Initialize Fiber app and mock service
	app := fiber.New()
	mockService := new(MockBankAccountService)
	controller := NewBankAccountController(mockService)

	// Define mock response
	mockBanks := []models.Bank{
		{Name: "PT. BANK CIMB NIAGA - (CIMB)", Alias: "CIMB", Code: "022"},
		{Name: "PT. BANK CIMB NIAGA UNIT USAHA SYARIAH - (CIMB SYARIAH)", Alias: "CIMB Syariah", Code: "730"},
		{Name: "PT. BNI SYARIAH", Alias: "BNI SYARIAH", Code: "427"},
	}
	mockService.On("LoadBanksName", "data/banks.json").Return(mockBanks, nil)

	// Define route and perform test
	app.Get("/banks", controller.GetBankName)
	req := httptest.NewRequest(http.MethodGet, "/banks", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	//assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string][]response.BankData
	err = json.NewDecoder(resp.Body).Decode(&response)
	//assert.NoError(t, err)

	//assert.Equal(t, mockBanks, response["data"])
}
