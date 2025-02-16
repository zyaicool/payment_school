package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"schoolPayment/controllers"
	request "schoolPayment/dtos/request"
	"schoolPayment/models"
	"schoolPayment/routes"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for InvoiceFormatService
type MockInvoiceFormatService struct {
	mock.Mock
}

func (m *MockInvoiceFormatService) Create(req *request.CreateInvoiceFormatRequest, userID int) (*models.InvoiceFormat, error) {
	args := m.Called(req, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InvoiceFormat), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockInvoiceFormatService) GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error) {
	args := m.Called(schoolID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.InvoiceFormat), args.Error(1)
	}
	return nil, args.Error(1)
}

func generateJWTInvoiceFormat(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

func TestCreateInvoiceFormat_Success(t *testing.T) {
    app := fiber.New()

    // Setup route
    routes.SetupInvoiceFormatRoutes(app, &controllers.InvoiceFormatController{})

    // Prepare request body
    requestBody := []byte(`{"schoolID": 1, "format": "some format"}`)

    claims := jwt.MapClaims{
        "user_id": 0,
    }
    token := generateJWTInvoiceFormat(claims)

    app.Use(func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak ditemukan")
        }
    
        // Verifikasi token JWT menggunakan library jwt-go
        tokenString := strings.TrimPrefix(token, "Bearer ")
        tokenClaims, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        if !tokenClaims.Valid {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        claims, ok := tokenClaims.Claims.(jwt.MapClaims)
        if !ok {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        c.Locals("user", claims)
        return c.Next()
    })

    req := httptest.NewRequest("POST", "/invoiceFormat/create", bytes.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    // Run request
    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("Test failed: %v", err)
    }

    // Check if the message returned matches the expected one
    var responseBody map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
        t.Fatalf("Failed to decode response body: %v", err)
    }
}

func TestCreateInvoiceFormat_Unauthorized(t *testing.T) {
	app := fiber.New()
	mockService := new(MockInvoiceFormatService)
	controller := controllers.NewInvoiceFormatController(mockService)

	app.Post("/invoiceFormat/create", controller.Create)

	// Mock data
	mockRequest := request.CreateInvoiceFormatRequest{
		SchoolID: 1,
		Prefix:   "INV",
		Format:   "YYYY-MM",
	}

	body, _ := json.Marshal(mockRequest)
	req := httptest.NewRequest("POST", "/invoiceFormat/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var responseBody map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Equal(t, "Unauthorized or invalid token", responseBody["error"])
}

//---cara kedua
func TestGetInvoiceFormatBySchoolID_Success(t *testing.T) {
    app := fiber.New()

    // Setup route
    routes.SetupInvoiceFormatRoutes(app, &controllers.InvoiceFormatController{})

    claims := jwt.MapClaims{
        "user_id": 0,
    }
    token := generateJWTInvoiceFormat(claims)

    app.Use(func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak ditemukan")
        }
    
        // Verifikasi token JWT menggunakan library jwt-go
        tokenString := strings.TrimPrefix(token, "Bearer ")
        tokenClaims, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        if !tokenClaims.Valid {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        claims, ok := tokenClaims.Claims.(jwt.MapClaims)
        if !ok {
            return fiber.NewError(fiber.StatusUnauthorized, "Token JWT tidak valid")
        }
    
        c.Locals("user", claims)
        return c.Next()
    })

    req := httptest.NewRequest("GET", "/invoiceFormat/detail", nil)
    req.Header.Set("Authorization", "Bearer "+token) // Set valid JWT token

    // Run request
    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("Test failed: %v", err)
    }

    // Check response body
    var responseBody map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&responseBody)
}

func TestGetInvoiceFormatBySchoolID_NotFound(t *testing.T) {
	app := fiber.New()
	mockService := new(MockInvoiceFormatService)
	controller := controllers.NewInvoiceFormatController(mockService)

	app.Get("/invoiceFormat/detail", controller.GetBySchoolID)

	// Mock data
	mockService.On("GetBySchoolID", uint(1)).Return(nil, errors.New("record not found"))

	req := httptest.NewRequest("GET", "/invoiceFormat/detail?schoolId=1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var responseBody map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.Equal(t, "record not found", responseBody["error"])
}
