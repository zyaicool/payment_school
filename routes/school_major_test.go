package routes

import (
	controllers "schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupSchoolMajorRoutes(t *testing.T) {
	// Inisialisasi Fiber app
	app := fiber.New()
	api := app.Group("/api/v1") 

	// Inisialisasi controller (bisa menggunakan controller kosong untuk testing)
	schoolMajorController := &controllers.SchoolMajorController{}

	// Setup routes
	SetupSchoolMajorRoutes(api, schoolMajorController)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Define expected routes with their methods and paths separately
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/schoolMajor/getListSchoolMajor"},
		{"POST", "/api/v1/schoolMajor/create"},           
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