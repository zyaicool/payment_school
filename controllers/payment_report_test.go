package controllers

import (
	"bytes"
	"math/big"
	"net/http/httptest"
	response "schoolPayment/dtos/response"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentReportService struct {
	mock.Mock
}

func (m *MockPaymentReportService) GetPaymentReport(page, limit int, sortBy, sortOrder string, paymentTypeId, userId int, startDate, endDate time.Time, studentId int, isAllData bool, userLoginID int, paymentStatus string, schoolGradeId int, schoolClassId, schoolLoginId int) (response.PaymentReportResponse, error) {
	args := m.Called(page, limit, sortBy, sortOrder, paymentTypeId, userId, startDate, endDate, studentId, isAllData, userLoginID, paymentStatus, schoolGradeId, schoolClassId, schoolLoginId)
	return args.Get(0).(response.PaymentReportResponse), args.Error(1)
}

func (m *MockPaymentReportService) ExportToExcel(paymentTypeId, userId int, startDate, endDate time.Time, studentId int, ctx *fiber.Ctx, paymentStatus string, schoolGradeId int, schoolClassId int) (*bytes.Buffer, error) {
	args := m.Called(paymentTypeId, userId, startDate, endDate, studentId, ctx, paymentStatus, schoolGradeId, schoolClassId) //(paymentTypeId, userId, studentId, paymentStatus, schoolGradeId, schoolClassId)
	return args.Get(0).(*bytes.Buffer), args.Error(1)
}

func parseTimeExcel(timeString string) *time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return nil
	}
	return &parsedTime
}

func TestGetPaymentReport(t *testing.T) {
	app := fiber.New()
	mockService := new(MockPaymentReportService)
	controller := NewPaymentReportController(mockService)

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{
			"user_id":   float64(1), // Mocked user ID
			"role_id":   float64(4),
			"school_id": float64(1), // Mocked school ID
		})
		return c.Next()
	})

	// Mock response
	mockReportResponse := response.PaymentReportResponse{
		TotalTransactionAmount: big.NewInt(1000000),
		TotalTransaction:       10,
		TotalStudent:           5,
		ListPaymentReport: response.ListPaymentReport{
			Page:      1,
			Limit:     10,
			TotalPage: 1,
			TotalData: 10,
			Data: []response.PaymentReportDetail{
				{
					ID:                1,
					InvoiceNumber:     "INV12345",
					StudentName:       "John Doe",
					PaymentDate:       nil,
					PaymentMethod:     "Bank Transfer",
					Username:          "admin",
					SchoolGradeName:   "Grade 10",
					SchoolClassName:   "Class A",
					TotalAmount:       500000,
					TransactionStatus: "Completed",
				},
			},
		},
	}

	mockService.On(
		"GetPaymentReport",
		1, 10, "", "asc", 0, 0, mock.Anything, mock.Anything,
		0, false, 1, "", 0, 0, 1,
	).Return(mockReportResponse, nil).Once()

	app.Get("/paymentReport/getList", controller.GetPaymentReport)

	req := httptest.NewRequest("GET", "/paymentReport/getList?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer mock-token")

	// Send the request and get the response
	resp, err := app.Test(req)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the response status is OK (200)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify that the mock method was called as expected
	mockService.AssertExpectations(t)
}

func TestExportToExcel(t *testing.T) {
	app := fiber.New()
	mockService := new(MockPaymentReportService)
	controller := NewPaymentReportController(mockService)

	// Middleware to set mock JWT claims
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{
			"user_id":   float64(1), // Mocked user ID
			"role_id":   float64(4),
			"school_id": float64(1), // Mocked school ID
		})
		return c.Next()
	})
	// Gunakan fungsi parseTime untuk mengonversi string ke time.Time
	startTime := parseTimeExcel("0001-01-01T00:00:00Z")
	endTime := parseTimeExcel("0001-01-01T00:00:00Z")

	if startTime == nil || endTime == nil {
		t.Fatalf("Failed to parse start or end time")
	}
	// Mock response
	mockBuffer := &bytes.Buffer{}
	mockBuffer.WriteString("mock excel content")

	// mockService.On(
	// 	"ExportToExcel",
	// 	0, 1, *startTime, *endTime, 0, mock.Anything, "", 0, 0,
	// ).Return(mockBuffer, nil).Once()
	mockService.On(
		"ExportToExcel",
		0,                                 // paymentTypeId
		0,                                 // userId
		*startTime,                        // startDate
		*endTime,                          // endDate
		0,                                 // studentId
		mock.AnythingOfType("*fiber.Ctx"), // Gunakan mock.AnythingOfType untuk ctx
		"",                                // paymentStatus
		0,                                 // schoolGradeId
		0,                                 // schoolClassId
	).Return(mockBuffer, nil).Once()

	app.Get("/paymentReport/exportExcel", controller.ExportPaymentReportToExcel)

	req := httptest.NewRequest("GET", "/paymentReport/exportExcel", nil)
	req.Header.Set("Authorization", "Bearer mock-token")

	// Send the request and get the response
	resp, err := app.Test(req)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the response status is OK (200)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	//Assert the response body contains mock excel content
	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	assert.Contains(t, body.String(), "mock excel content")

	// Verify that the mock method was called as expected
	mockService.AssertExpectations(t)
}
