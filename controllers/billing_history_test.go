package controllers

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"schoolPayment/dtos/response"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// First, ensure the mock implements the interface
type MockBillingHistoryService struct {
	mock.Mock
}

// Add this line to verify interface implementation at compile time
var _ services.BillingHistoryServiceInterface = (*MockBillingHistoryService)(nil)

// Implement all required methods from the interface
func (m *MockBillingHistoryService) GetAllBillingHistory(
	page int,
	limit int,
	search string,
	studentID int,
	roleID int,
	schoolYearId int,
	paymentTypeId int,
	paymentStatusCode string,
	userID int,
	userLoginID int,
	sortBy string,
	sortOrder string,
) (response.ListBillingHistory, error) {
	args := m.Called(page, limit, search, studentID, roleID, schoolYearId, paymentTypeId, paymentStatusCode, userID, userLoginID, sortBy, sortOrder)
	return args.Get(0).(response.ListBillingHistory), args.Error(1)
}

func (m *MockBillingHistoryService) GetDetailBillingHistoryIDService(
	transactionId int,
	userID int,
) (*response.DetailBillingHistory, error) {
	args := m.Called(transactionId, userID)
	return args.Get(0).(*response.DetailBillingHistory), args.Error(1)
}

func (m *MockBillingHistoryService) GenerateInvoice(
	c *fiber.Ctx,
	invoiceNumber string,
	isPrint bool,
	userID int,
) (string, error) {
	args := m.Called(c, invoiceNumber, isPrint, userID)
	return args.String(0), args.Error(1)
}
func TestGetAllBillingHistory(t *testing.T) {
	// Create mock service
	mockService := new(MockBillingHistoryService)

	// Create controller with mock service
	controller := NewBillingHistoryController(mockService)

	// Initialize Fiber app
	app := fiber.New()

	// IMPORTANT: Add middleware before defining routes
	app.Use(func(c *fiber.Ctx) error {
		// Mock JWT claims
		claims := jwt.MapClaims{
			"user_id": float64(123), // must be float64 to match jwt.MapClaims type
			"role_id": float64(1),   // must be float64 to match jwt.MapClaims type
		}
		c.Locals("user", claims)
		return c.Next()
	})

	// Add route after middleware
	app.Get("/api/v1/billingHistory/getAllBillingHistory", controller.GetAllBillingHistory)

	// Mock expected service response
	now := time.Now()
	expectedResponse := response.ListBillingHistory{
		Page:      1,
		Limit:     10,
		TotalData: int64(100),
		TotalPage: 10,
		Data: []response.DataListBillingHistory{
			{
				ID:                1,
				InvoiceNumber:     "INV123",
				StudentName:       "John Doe",
				PaymentDate:       &now,
				PaymentMethod:     "kasir",
				Username:          "admin",
				TotalAmount:       1000,
				TransactionStatus: "menunggu",
				OrderID:           "ORD123",
				RedirectUrl:       "http://example.com",
				Token:             "token123",
			},
		},
	}

	// Setup mock service expectations with exact parameters
	mockService.On("GetAllBillingHistory",
		1,     // page
		10,    // limit
		"",    // search
		0,     // studentID
		1,     // roleID - from mock JWT claims
		0,     // schoolYearId
		0,     // paymentTypeId
		"",    // paymentStatusCode
		0,     // userID
		123,   // userLoginID - from mock JWT claims
		"",    // sortBy
		"asc", // sortOrder
	).Return(expectedResponse, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/billingHistory/getAllBillingHistory?page=1&limit=10", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response body
	var actualResponse response.ListBillingHistory
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)

	// Response assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Page, actualResponse.Page)
	assert.Equal(t, expectedResponse.Limit, actualResponse.Limit)
	assert.Equal(t, expectedResponse.TotalData, actualResponse.TotalData)
	assert.Equal(t, expectedResponse.TotalPage, actualResponse.TotalPage)
	assert.Len(t, actualResponse.Data, 1)

	// Verify the first data item
	actualData := actualResponse.Data[0]
	expectedData := expectedResponse.Data[0]
	assert.Equal(t, expectedData.ID, actualData.ID)
	assert.Equal(t, expectedData.InvoiceNumber, actualData.InvoiceNumber)
	assert.Equal(t, expectedData.StudentName, actualData.StudentName)
	assert.Equal(t, expectedData.PaymentMethod, actualData.PaymentMethod)
	assert.Equal(t, expectedData.Username, actualData.Username)
	assert.Equal(t, expectedData.TotalAmount, actualData.TotalAmount)
	assert.Equal(t, expectedData.TransactionStatus, actualData.TransactionStatus)
	assert.Equal(t, expectedData.OrderID, actualData.OrderID)
	assert.Equal(t, expectedData.RedirectUrl, actualData.RedirectUrl)
	assert.Equal(t, expectedData.Token, actualData.Token)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}
func TestGetDetailBillingHistoryID(t *testing.T) {
	// Create mock service
	mockService := new(MockBillingHistoryService)

	// Create controller with mock service
	controller := NewBillingHistoryController(mockService)

	// Initialize Fiber app
	app := fiber.New()

	// Add middleware before routes
	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": float64(123),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	// Add route after middleware
	app.Get("/api/v1/billingHistory/detailBillingHistory", controller.GetDetailBillingHistoryID)

	// Mock expected service response
	now := time.Now()
	expectedResponse := &response.DetailBillingHistory{
		ID:                 1,
		StudentName:        "John Doe",
		SchoolClass:        "Class A",
		InvoiceNumber:      "INV123",
		TransactionStatus:  "completed",
		ChangeAmount:       0,
		TotalBillingAmount: 1000,
		TotalPayAmount:     1000,
		PaymentDate:        &now,
		AdminFee:           10,
		PaymentMethod:      "Cash",
		ListBilling:        []response.BillingStudentForHistory{},
	}

	// Setup mock service expectations
	mockService.On("GetDetailBillingHistoryIDService", 1, 123).Return(expectedResponse, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/billingHistory/detailBillingHistory?transactionId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Parse response body
	var actualResponse response.DetailBillingHistory
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)

	// Response assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, actualResponse.ID)
	assert.Equal(t, expectedResponse.StudentName, actualResponse.StudentName)
	assert.Equal(t, expectedResponse.SchoolClass, actualResponse.SchoolClass)
	assert.Equal(t, expectedResponse.InvoiceNumber, actualResponse.InvoiceNumber)
	assert.Equal(t, expectedResponse.TransactionStatus, actualResponse.TransactionStatus)
	assert.Equal(t, expectedResponse.TotalBillingAmount, actualResponse.TotalBillingAmount)
	assert.Equal(t, expectedResponse.TotalPayAmount, actualResponse.TotalPayAmount)
	assert.Equal(t, expectedResponse.AdminFee, actualResponse.AdminFee)
	assert.Equal(t, expectedResponse.PaymentMethod, actualResponse.PaymentMethod)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetDetailBillingHistoryID_NotFound(t *testing.T) {
	// Create mock service
	mockService := new(MockBillingHistoryService)

	// Create controller with mock service
	controller := NewBillingHistoryController(mockService)

	// Initialize Fiber app
	app := fiber.New()

	// Add middleware before routes
	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": float64(123),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	// Add route after middleware
	app.Get("/api/v1/billingHistory/detailBillingHistory", controller.GetDetailBillingHistoryID)

	// Setup mock service expectations with error
	mockService.On("GetDetailBillingHistoryIDService", 1, 123).Return((*response.DetailBillingHistory)(nil), fmt.Errorf("not found"))

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/billingHistory/detailBillingHistory?transactionId=1", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGeneratePDF(t *testing.T) {
	// Create mock service
	mockService := new(MockBillingHistoryService)

	// Create controller with mock service
	controller := NewBillingHistoryController(mockService)

	// Initialize Fiber app
	app := fiber.New()

	// Add middleware before routes
	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": float64(123),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	// Add route after middleware
	app.Get("/api/v1/billingHistory/printInvoice", controller.GeneratePDF)

	// Setup mock service expectations
	mockService.On("GenerateInvoice", mock.Anything, "INV123", true, 123).Return("invoice.pdf", nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/billingHistory/printInvoice?invoiceNumber=INV123", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGeneratePDF_Error(t *testing.T) {
	// Create mock service
	mockService := new(MockBillingHistoryService)

	// Create controller with mock service
	controller := NewBillingHistoryController(mockService)

	// Initialize Fiber app
	app := fiber.New()

	// Add middleware before routes
	app.Use(func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"user_id": float64(123),
		}
		c.Locals("user", claims)
		return c.Next()
	})

	// Add route after middleware
	app.Get("/api/v1/billingHistory/printInvoice", controller.GeneratePDF)

	// Setup mock service expectations with error
	mockService.On("GenerateInvoice", mock.Anything, "INV123", true, 123).Return("", fmt.Errorf("failed to generate PDF"))

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/billingHistory/printInvoice?invoiceNumber=INV123", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Parse response body
	var errorResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "failed to generate PDF", errorResponse["error"])

	// Verify mock expectations
	mockService.AssertExpectations(t)
}
