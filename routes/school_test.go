package routes

import (
	"schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupSchoolRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api/v1")
	mockController := &controllers.SchoolController{}

	// Setup the routes
	SetupSchoolRoutes(api, mockController)

	// Get the route stack
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/school/getAllSchool"},
		{"GET", "/api/v1/school/detail/:id"},
		{"POST", "/api/v1/school/create"},
		{"PUT", "/api/v1/school/update/:id"},
		{"DELETE", "/api/v1/school/delete/:id"},
		{"GET", "/api/v1/school/getAllOnboarding"},
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
