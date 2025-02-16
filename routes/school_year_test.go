package routes

import (
	controllers "schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupSchoolYearRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api/v1")
	mockController := &controllers.SchoolYearController{}

	SetupSchoolYearRoutes(api, mockController)

	stack := app.Stack()
	assert.NotEmpty(t, stack)

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/schoolYear/getAllSchoolYear"},
		{"GET", "/api/v1/schoolYear/detail/:id"},
		{"POST", "/api/v1/schoolYear/create"},
		{"PUT", "/api/v1/schoolYear/update/:id"},
		{"DELETE", "/api/v1/schoolYear/delete/:id"},
	}

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
