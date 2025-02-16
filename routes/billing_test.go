package routes

import (
	controllers "schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupBillingRoutes(t *testing.T) {
	// Inisialisasi Fiber app
	app := fiber.New()
	api := app.Group("/api/v1")

	
	billingController := &controllers.BillingController{}

	// Setup routes
	SetupBillingRoutes(api, billingController)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes with their methods and paths separately
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/billing/getAllBilling"},
		{"GET", "/api/v1/billing/detail/:id"},
		{"POST", "/api/v1/billing/create"},
		{"PUT", "/api/v1/billing/update/:id"},
		{"DELETE", "/api/v1/billing/delete/:id"},
		{"GET", "/api/v1/billing/generateInstallment/:id"},
		{"GET", "/api/v1/billing/billingStatus"},
		{"POST", "/api/v1/billing/createDonation"},
		{"GET", "/api/v1/billing/getBillingByStudentID"},
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