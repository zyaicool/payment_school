package controllers_test

import (
	"bytes"
	"mime/multipart"
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

type MockPaymentMethodService struct {
	mock.Mock
}

func (m *MockPaymentMethodService) CreatePaymentMethod(paymentRequest *request.PaymentMethodCreateRequest, userID int) (*models.PaymentMethod, error) {
	args := m.Called(paymentRequest, userID)
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodService) UpdatePaymentMethod(paymentMethodID int, updateRequest *request.PaymentMethodCreateRequest) (*models.PaymentMethod, error) {
	args := m.Called(paymentMethodID, updateRequest)
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodService) GetAllPaymentMethod(search string) ([]response.PaymentMethodResponse, error) {
	args := m.Called(search)
	return args.Get(0).([]response.PaymentMethodResponse), args.Error(1)
}

func (m *MockPaymentMethodService) GetPaymentMethodDetail(id int) (*response.PaymentMethodResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*response.PaymentMethodResponse), args.Error(1)
}

func TestCreatePaymentMethod(t *testing.T) {
	mockService := new(MockPaymentMethodService)

	// Mock Data
	mockPaymentMethod := &models.PaymentMethod{
		ID:                 1,
		PaymentMethod:      "Credit Card",
		BankCode:           "CC",
		BankName:           "Bank ABC",
		AdminFee:           1000,
		MethodLogo:         "logo.png",
		IsPercentage:       false,
		AdminFeePercentage: "",
	}
	mockService.On("CreatePaymentMethod", mock.Anything, mock.Anything).Return(mockPaymentMethod, nil)

	controller := controllers.NewPaymentMethodController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, 
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/masterPaymentGateway/create", controller.CreatePaymentMethod)

	// Create form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("paymentMethod", "Credit Card")
	writer.WriteField("bankCode", "CC")
	writer.WriteField("bankName", "Bank ABC")
	writer.WriteField("adminFee", "1000")
	writer.WriteField("isPercentage", "false")
	writer.WriteField("adminFeePercentage", "")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/masterPaymentGateway/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestUpdatePaymentMethod(t *testing.T) {
	mockService := new(MockPaymentMethodService)

	// Mock Data
	mockPaymentMethod := &models.PaymentMethod{
		ID:                 1,
		PaymentMethod:      "Updated Credit Card",
		BankCode:           "CC",
		BankName:           "Bank ABC",
		AdminFee:           1500,
		MethodLogo:         "updated_logo.png",
		IsPercentage:       false,
		AdminFeePercentage: "",
	}
	mockService.On("UpdatePaymentMethod", 1, mock.Anything).Return(mockPaymentMethod, nil)

	controller := controllers.NewPaymentMethodController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Put("/api/v1/masterPaymentGateway/update/:id", controller.UpdatePaymentMethod)

	// Create form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("paymentMethod", "Updated Credit Card")
	writer.WriteField("bankCode", "CC")
	writer.WriteField("bankName", "Bank ABC")
	writer.WriteField("adminFee", "1500")
	writer.WriteField("isPercentage", "false")
	writer.WriteField("adminFeePercentage", "")
	writer.Close()

	req := httptest.NewRequest("PUT", "/api/v1/masterPaymentGateway/update/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetAllPaymentMethod(t *testing.T) {
	mockService := new(MockPaymentMethodService)

	// Mock Data
	mockResponse := []response.PaymentMethodResponse{
		{
			ID:                 1,
			PaymentMethod:      "Credit Card",
			BankCode:           "CC",
			BankName:           "Bank ABC",
			AdminFee:           1000,
			MethodLogo:         "logo.png",
			IsPercentage:       false,
			AdminFeePercentage: "",
		},
	}
	mockService.On("GetAllPaymentMethod", "").Return(mockResponse, nil)

	controller := controllers.NewPaymentMethodController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/masterPaymentGateway/getAllConfig", controller.GetAllPaymentMethod)

	req := httptest.NewRequest("GET", "/api/v1/masterPaymentGateway/getAllConfig", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetPaymentMethodDetail(t *testing.T) {
	mockService := new(MockPaymentMethodService)

	// Mock Data
	mockResponse := &response.PaymentMethodResponse{
		ID:                 1,
		PaymentMethod:      "Credit Card",
		BankCode:           "CC",
		BankName:           "Bank ABC",
		AdminFee:           1000,
		MethodLogo:         "logo.png",
		IsPercentage:       false,
		AdminFeePercentage: "",
	}
	mockService.On("GetPaymentMethodDetail", 1).Return(mockResponse, nil)

	controller := controllers.NewPaymentMethodController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/masterPaymentGateway/detail/:id", controller.GetPaymentMethodDetail)

	req := httptest.NewRequest("GET", "/api/v1/masterPaymentGateway/detail/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func generateJWTPaymentMethod(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}