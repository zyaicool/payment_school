package routes

import (
	"testing"

	controllers "schoolPayment/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupTransactionRoutes(t *testing.T) {
	// Inisialisasi Fiber app
	app := fiber.New()
	api := app.Group("/api/v1") 

	// Inisialisasi controller (bisa menggunakan controller kosong untuk testing)
	transactionController := &controllers.TransactionController{}

	// Setup routes
	SetupTransactionRoutes(api, transactionController)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes with their methods and paths separately
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/transaction/create"},               
		{"POST", "/api/v1/transaction/donation"},             
		{"POST", "/api/v1/transaction/midtrans/payment"},    
		{"GET", "/api/v1/transaction/midtrans/checkPayment"},
		{"POST", "/api/v1/webhook"},                         
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