package routes

import (
	controllers "schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupBankAccountRoutes(t *testing.T) {
	// Inisialisasi Fiber app
	app := fiber.New()
	api := app.Group("/api/v1") 


	bankAccountController := &controllers.BankAccountController{}

	// Setup routes
	SetupBankAccountRoutes(api, bankAccountController)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes with their methods and paths
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/bankAccount/getListBankAccount"}, 
		{"GET", "/api/v1/bankAccount/detail/:id"},        
		{"POST", "/api/v1/bankAccount/create"},          
		{"PUT", "/api/v1/bankAccount/update/:id"},         
		{"DELETE", "/api/v1/bankAccount/delete/:id"},      
		{"GET", "/api/v1/bankAccount/listBankName"},      
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