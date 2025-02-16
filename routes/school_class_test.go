package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupSchoolClassRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api/v1")

	SetupSchoolClassRoutes(api)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes with their methods and paths separately
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/schoolClass/getAllSchoolClass"},
		{"GET", "/api/v1/schoolClass/detail/:id"},
		{"POST", "/api/v1/schoolClass/create"},
		{"PUT", "/api/v1/schoolClass/update/:id"},
		{"DELETE", "/api/v1/schoolClass/delete/:id"},
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
