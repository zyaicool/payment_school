package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"schoolPayment/dtos/response"
	"schoolPayment/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllRoles(page int, limit int, search string, roleID int) (response.RoleListResponse, error) {
	args := m.Called(page, limit, search, roleID)
	return args.Get(0).(response.RoleListResponse), args.Error(1)
}

func TestGetAllRole(t *testing.T) {
	app := fiber.New()

	mockService := new(MockService)

	expectedRoles := response.RoleListResponse{
		Limit: 10,
		Page:  1,
		Data: []models.Role{
			{Master: models.Master{ID: 1}, Name: "Admin"},
			{Master: models.Master{ID: 2}, Name: "User"},
		},
	}

	userDataClaims := jwt.MapClaims{
		"user_id": float64(1),
		"role_id": float64(1),
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", userDataClaims)
		return c.Next()
	})

	// Create a route for GetAllRole
	app.Get("/v1/roles/getAllRole", func(c *fiber.Ctx) error {
		userClaims := c.Locals("user").(jwt.MapClaims)
		roleID := int(userClaims["role_id"].(float64))

		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		search := c.Query("search")

		roles, err := mockService.GetAllRoles(page, limit, search, roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot fetch data",
			})
		}
		return c.JSON(roles)
	})

	t.Run("Successful Retrieval of Roles", func(t *testing.T) {
		// Set up the mock to return expected roles
		mockService.On("GetAllRoles", 1, 10, "", 1).Return(expectedRoles, nil)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, "/v1/roles/getAllRole?page=1&limit=10", nil)

		// Set the user role in the context
		userClaims := jwt.MapClaims{
			"role_id": float64(1), // Assuming role_id is 1
		}
		req = req.WithContext(context.WithValue(req.Context(), "user", userClaims))

		// Send the request to the app
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		// Check the status code
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check the response body for the roles data
		var rolesResponse response.RoleListResponse
		err = json.NewDecoder(resp.Body).Decode(&rolesResponse)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, expectedRoles, rolesResponse)

		// Assert that the mock service was called as expected
		mockService.AssertCalled(t, "GetAllRoles", 1, 10, "", 1)
	})
}

func (m *MockService) GetRoleByID(id uint) (models.Role, error) {
	args := m.Called(id)
	return args.Get(0).(models.Role), args.Error(1)
}

func TestGetDataRole(t *testing.T) {
	// Create a new Fiber instance
	app := fiber.New()

	// Create a mock service
	mockService := new(MockService)

	// Set up the route with the handler function
	app.Get("/v1/role/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		role, err := mockService.GetRoleByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data not found",
			})
		}
		return c.JSON(role)
	})

	t.Run("Role Found", func(t *testing.T) {
		// Define the mock response data
		expectedRole := models.Role{
			Master: models.Master{ID: 1},
			Name:   "Admin",
		}

		// Set up the mock to return the expected role when GetRoleByID is called
		mockService.On("GetRoleByID", uint(1)).Return(expectedRole, nil)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, "/v1/role/1", nil)

		// Send the request to the app
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		// Check the status code
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Check the response body
		var roleResponse models.Role
		err = json.NewDecoder(resp.Body).Decode(&roleResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedRole, roleResponse)

		// Assert that the mock service was called as expected
		mockService.AssertCalled(t, "GetRoleByID", uint(1))
	})

	t.Run("Role Not Found", func(t *testing.T) {
		// Set up the mock to return an error for a non-existent role
		mockService.On("GetRoleByID", uint(2)).Return(models.Role{}, fmt.Errorf("Data not found"))

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, "/v1/role/2", nil)

		// Send the request to the app
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		// Check the status code
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		// Check the response body
		var errorResponse map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Data not found", errorResponse["error"])

		// Assert that the mock service was called as expected
		mockService.AssertCalled(t, "GetRoleByID", uint(2))
	})
}
