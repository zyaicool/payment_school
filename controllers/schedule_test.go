package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"schoolPayment/controllers"
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockScheduleService struct {
	mock.Mock
}

func (m *MockScheduleService) GetCheckPaymentFailedUsingScheduleService() (response.CheckingPaymentStatusResponse, error) {
	args := m.Called()
	return args.Get(0).(response.CheckingPaymentStatusResponse), args.Error(1)
}

func (m *MockScheduleService) DummySendNotif(userId int, schedule *request.DummyNotifRequest) error {
	args := m.Called(userId, schedule)
	return args.Error(0)
}

func (m *MockScheduleService) SendReminderDueDate() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetCheckPaymentFailedUsingSchedule(t *testing.T) {
	mockService := new(MockScheduleService)

	// Mock Data
	mockData := response.CheckingPaymentStatusResponse{
		Data: []response.CheckingPaymentStatusDetailResponse{
			{
				ID:                   1,
				TransactionType:      "Type 1",
				VirtualAccountNumber: "1234567890",
				TotalAmount:          10000,
				ReferenceNumber:      "REF123",
				Description:          "Test Transaction",
				StudentID:            1,
				OrderID:              "ORDER123",
				TransactionStatus:    "PS01",
				InvoiceNumber:        "INV123",
				BillingStudentIds:    "1,2,3",
			},
		},
	}
	mockService.On("GetCheckPaymentFailedUsingScheduleService").Return(mockData, nil)

	controller := controllers.NewScheduleController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/schedule/checkFailedTransaction", controller.GetCheckPaymentFailedUsingSchedule)

	req := httptest.NewRequest("GET", "/api/v1/schedule/checkFailedTransaction", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSendDummyNotif(t *testing.T) {
	mockService := new(MockScheduleService)

	// Mock Data
	mockRequest := &request.DummyNotifRequest{
		Title: "Test Title",
		Body:  "Test Body",
	}
	mockService.On("DummySendNotif", 5, mockRequest).Return(nil)

	controller := controllers.NewScheduleController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Post("/api/v1/schedule/sendDummyNotif", controller.SendDummyNotif)

	body, _ := json.Marshal(mockRequest)
	req := httptest.NewRequest("POST", "/api/v1/schedule/sendDummyNotif", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestSendBillingReminder(t *testing.T) {
	mockService := new(MockScheduleService)

	// Mock Data
	mockService.On("SendReminderDueDate").Return("Email Send", nil)

	controller := controllers.NewScheduleController(mockService)

	app := fiber.New()

	claims := jwt.MapClaims{
		"user_id": 1.0, // Sesuaikan dengan user ID yang diharapkan
	}
	token := generateJWT(claims)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", claims)
		return c.Next()
	})

	app.Get("/api/v1/schedule/sendBillingReminder", controller.SendBillingReminder)

	req := httptest.NewRequest("GET", "/api/v1/schedule/sendBillingReminder", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func generateJWTSchedule(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

