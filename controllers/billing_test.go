package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking BillingService
type MockBillingService struct {
	mock.Mock
}

func (m *MockBillingService) GetBillingByID(id uint) (response.BillingDetailResponse, error) {
	args := m.Called(id)
	return args.Get(0).(response.BillingDetailResponse), args.Error(1)
}

func (m *MockBillingService) GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, userID int) (response.BillingListResponse, error) {
	args := m.Called(page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, isDonation, userID)
	return args.Get(0).(response.BillingListResponse), args.Error(1)
}

func (m *MockBillingService) CreateBilling(billingRequest *request.BillingCreateRequest, userID int) (*models.Billing, error) {
	args := m.Called(billingRequest, userID)
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingService) GetBillingStatuses(filePath string) ([]models.BillingStatus, error) {
	args := m.Called(filePath)
	return args.Get(0).([]models.BillingStatus), args.Error(1)
}

func (m *MockBillingService) CreateDonation(billingName string, schoolGradeId int, bankAccountId int, userId uint) (response.BillingResponse, error) {
	args := m.Called(billingName, schoolGradeId, bankAccountId, userId)
	return args.Get(0).(response.BillingResponse), args.Error(1)
}

func (m *MockBillingService) GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) (*response.BillingByStudentResponse, error) {
	args := m.Called(studentID, schoolYearID, schoolGradeID, schoolClassID)
	return args.Get(0).(*response.BillingByStudentResponse), args.Error(1)
}

// Unit test for GetDataBilling controller
func TestGetDataBilling(t *testing.T) {
	// Setup Fiber app and controller
	app := fiber.New()
	mockService := new(MockBillingService)

	// Mock behavior of GetBillingByID
	mockService.On("GetBillingByID", uint(1)).Return(response.BillingDetailResponse{
		ID:              1,
		BillingName:     "Tuition Fee",
		BillingCode:     "TF123",
		BillingType:     "Tuition",
		BankAccountName: "Bank ABC - 123456",
		Description:     "Description of the fee",
		SchoolYear:      "2024",
		SchoolClassList: "Class 1A, Class 2B",
		DetailBillings: []response.DetailBilling{
			{ID: 1, DetailBillingName: "Tuition Fee", DueDate: time.Now(), Amount: 1000},
		},
	}, nil)

	// Test success scenario
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/billing/detail/1", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp, err := app.Test(req)

		assert.Nil(t, err)

		var body response.BillingDetailResponse
		_ = json.NewDecoder(resp.Body).Decode(&body)
	})

	// Test Unauthorized access (if CheckAccessUserTuKasirAdminSekolah fails)
	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/billing/detail/1", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		resp, err := app.Test(req)

		assert.Nil(t, err)
		

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
	})


	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetBillingByID", uint(1)).Return(response.BillingDetailResponse{}, fmt.Errorf("data not found"))

		req := httptest.NewRequest("GET", "/api/v1/billing/detail/1", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, 404, resp.StatusCode)

		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
	})
}

func TestCreateBilling_Success(t *testing.T) {
	// Initialize the app
	app := fiber.New()

	// Mock service
	mockService := new(MockBillingService)
	controller := NewBillingController(mockService)

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{
			"user_id": float64(1),  // Mocked user ID
			"role_id": float64(4), // Mocked role ID (Kasir, which should pass)
		})
		return c.Next()
	})

	// Register the route
	app.Post("/create", controller.CreateBilling)

	// Mock the service method
	mockService.On("CreateBilling", mock.AnythingOfType("*request.BillingCreateRequest"), 1).
		Return(&models.Billing{
			BillingName: "Tuition Fee",
		}, nil)

	// Prepare request body
	billingRequest := request.BillingCreateRequest{
		BillingType:    "Tuition",
		SchoolGradeID:  1,
		SchoolYearId:   2023,
		BillingName:    "Tuition Fee",
		BillingAmount:  1000,
		Description:    "Monthly tuition fee",
		BillingCode:    "TUI2023",
		SchoolClassIds: []string{"1", "2"},
		BankAccountId:  1,
		DetailBillings: []request.DetailBillings{},
	}
	requestBody, _ := json.Marshal(billingRequest)

	req := httptest.NewRequest("POST", "/create", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the response message and data
	assert.Equal(t, "Data berhasil disimpan.", responseBody["message"])
	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, "Tuition Fee", data["billingName"])

	// Assert the mock expectations
	mockService.AssertExpectations(t)
}

func TestGetAllBilling(t *testing.T) {
	// Initialize Fiber app
	app := fiber.New()

	// Mock dependencies
	mockService := new(MockBillingService)
	controller := NewBillingController(mockService)

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{
			"user_id": float64(1),  // Mocked user ID
		})
		return c.Next()
	})

	// Define the route
	app.Get("/billing", controller.GetAllBilling)

	// Define test data
	page := 1
	limit := 10
	search := "Test"
	billingType := "Type1"
	paymentType := "Online"
	schoolGrade := "1"
	bankAccountId := 1
	sort := "desc"
	sortBy := "created_at"
	sortOrder := "asc"
	isDonation := false
	userID := 1

	mockResponse := response.BillingListResponse{
		Page:       page,
		Limit:      limit,
		TotalData:  2,
		TotalPage:  1,
		Data: []response.BillingResponse{
			{
				ID:              1,
				BillingName:     "Test Billing 1",
				BillingType:     "Type1",
				BankAccountName: "Bank 1 - 12345678",
				SchoolGradeName: "Grade 1",
				CreatedBy:       "Admin",
			},
			{
				ID:              2,
				BillingName:     "Test Billing 2",
				BillingType:     "Type2",
				BankAccountName: "Bank 2 - 87654321",
				SchoolGradeName: "Grade 2",
				CreatedBy:       "Admin",
			},
		},
	}

	// Mock the service response
	mockService.On("GetAllBilling", page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, &isDonation, userID).
		Return(mockResponse, nil)

	// Create a test request
	req := httptest.NewRequest("GET", "/billing?page=1&limit=10&search=Test&billingType=Type1&paymentType=Online&schoolGrade=1&bankAccountId=1&sort=desc&sortBy=created_at&sortOrder=asc&isDonation=false", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken") // Simulate JWT if required

	// Perform the request
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Parse response body
	var actualResponse response.BillingListResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, mockResponse, actualResponse)

	// Assert expectations
	mockService.AssertExpectations(t)
}