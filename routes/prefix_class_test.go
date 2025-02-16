package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPrefixClassController is a mock of the PrefixClassController
type MockPrefixClassController struct {
	mock.Mock
}

func (m *MockPrefixClassController) GetPrefixClass(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockPrefixClassController) CreatePrefixClass(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestSetupPrefixClassRoutes(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create API group
	api := app.Group("/api")

	// Test route setup
	t.Run("Routes are set up correctly", func(t *testing.T) {
		// Setup routes
		SetupPrefixClassRoutes(api)

		// Get the stack of registered routes
		routes := app.GetRoutes()

		// Verify routes are registered
		var foundGetRoute, foundPostRoute bool
		expectedGetPath := "/api/prefixClass/getListPrefixClass"
		expectedPostPath := "/api/prefixClass/create"

		for _, route := range routes {
			switch {
			case route.Path == expectedGetPath && route.Method == "GET":
				foundGetRoute = true
			case route.Path == expectedPostPath && route.Method == "POST":
				foundPostRoute = true
			}
		}

		// Assert that both routes were found
		assert.True(t, foundGetRoute, "GET route should be registered")
		assert.True(t, foundPostRoute, "POST route should be registered")
	})

	// Test middleware presence
	t.Run("Routes have JWT middleware", func(t *testing.T) {
		// Setup routes
		SetupPrefixClassRoutes(api)

		// Get routes
		routes := app.GetRoutes()

		// Check for middleware in routes
		for _, route := range routes {
			if route.Path == "/api/prefixClass/getListPrefixClass" ||
				route.Path == "/api/prefixClass/create" {
				// Verify that the route has handlers (middleware + endpoint)
				assert.GreaterOrEqual(t, len(route.Handlers), 2,
					"Route should have JWT middleware and handler")
			}
		}
	})
}
