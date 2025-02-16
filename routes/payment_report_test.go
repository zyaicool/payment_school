package routes

import (
	"schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupPaymentReportRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api/v1")
	mockController := &controllers.PaymentReportController{}

	// Setup the routes
	SetupPaymentReportRoutes(api, mockController)

	// Get the route stack
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/paymentReport/getList"},
		{"GET", "/api/v1/paymentReport/exportExcel"},
	}

	// Verify each expected route
	for _, expectedRoute := range expectedRoutes {
		found := false
		for _, routeStack := range stack {
			for _, route := range routeStack {
				if route.Method == expectedRoute.method && route.Path == expectedRoute.path {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		assert.True(t, found, "Route %s %s should be registered",
			expectedRoute.method, expectedRoute.path)
	}
}
