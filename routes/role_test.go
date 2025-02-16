package routes

import (
	controllers "schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoleRoutes(t *testing.T) {
	// Initialize Fiber app and mock controller
	app := fiber.New()
	api := app.Group("/api/v1")
	mockController := &controllers.RoleController{}

	// Setup routes with mock controller
	SetupRoleRoutes(api, mockController)

	// Get the route stack
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/role/getAllRole"},
		{"GET", "/api/v1/role/detail/:id"},
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
