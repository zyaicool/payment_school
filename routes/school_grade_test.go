package routes

import (
	"schoolPayment/controllers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSetupSchoolGradeRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api/v1")
	mockController := &controllers.SchoolGradeController{}

	SetupSchoolGradeRoutes(api, mockController)

	stack := app.Stack()
	assert.NotEmpty(t, stack)

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/schoolGrade/getAllSchoolGrade"},
		{"GET", "/api/v1/schoolGrade/detail/:id"},
		{"POST", "/api/v1/schoolGrade/create"},
		{"PUT", "/api/v1/schoolGrade/update/:id"},
		{"DELETE", "/api/v1/schoolGrade/delete/:id"},
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
