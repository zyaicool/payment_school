package services

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"schoolPayment/dtos/response"
// 	"schoolPayment/models"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockService struct {
// 	mock.Mock
// }

// func (m *MockService) GetAllUser(page int, limit int, search string, roleID int, userID int, schoolID int, sort string) (response.UserListResponse, error) {
// 	args := m.Called(page, limit, search, roleID, userID, schoolID, sort)
// 	return args.Get(0).(response.UserListResponse), args.Error(1)
// }

// func TestGetAllUser(t *testing.T) {
// 	app := fiber.New()

// 	mockService := new(MockService)

// 	expectedUser := response.UserListResponse{
// 		Limit:     10,
// 		Page:      1,
// 		TotalPage: 1,
// 		TotalData: 2,
// 		Data: []response.ListDataUserForIndex{
// 			{ID: 1, Username: "Admin", RoleID: 1, Email: "admin@example.com", Status: "Active"},
// 			{ID: 2, Username: "User", RoleID: 2, Email: "user@example.com", Status: "Active"},
// 		},
// 	}

// 	userDataClaims := jwt.MapClaims{
// 		"user_id": float64(1),
// 		"role_id": float64(1),
// 	}

// 	app.Use(func(c *fiber.Ctx) error {
// 		c.Locals("user", userDataClaims)
// 		return c.Next()
// 	})

// 	// Create a route for GetAllUser
// 	app.Get("/v1/user/getAllUser", func(c *fiber.Ctx) error {
// 		userClaims := c.Locals("user").(jwt.MapClaims)
// 		roleID := int(userClaims["role_id"].(float64))

// 		page := c.QueryInt("page", 1)
// 		limit := c.QueryInt("limit", 10)
// 		search := c.Query("search", "")
// 		schoolID := c.QueryInt("schoolId", 0)
// 		sort := c.Query("sort", "desc")

// 		users, err := mockService.GetAllUser(page, limit, search, roleID, 1, schoolID, sort)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Cannot fetch data",
// 			})
// 		}
// 		return c.JSON(users)
// 	})

// 	t.Run("Successful Retrieval of Users", func(t *testing.T) {
// 		// Set up the mock to return expected roles
// 		mockService.On("GetAllUser", 1, 10, "", 1, 1, 0, "desc").Return(expectedUser, nil)

// 		// Create a new HTTP request
// 		req := httptest.NewRequest(http.MethodGet, "/v1/user/getAllUser?page=1&limit=10", nil)

// 		// Set the user role in the context
// 		userClaims := jwt.MapClaims{
// 			"role_id": float64(1), // Assuming role_id is 1
// 		}
// 		req = req.WithContext(context.WithValue(req.Context(), "user", userClaims))

// 		// Send the request to the app
// 		resp, err := app.Test(req, -1)
// 		if err != nil {
// 			t.Fatalf("Failed to send request: %v", err)
// 		}

// 		// Check the status code
// 		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

// 		// Check the response body for the users data
// 		var usersResponse response.UserListResponse
// 		err = json.NewDecoder(resp.Body).Decode(&usersResponse)
// 		if err != nil {
// 			t.Fatalf("Failed to decode response: %v", err)
// 		}
// 		assert.Equal(t, expectedUser, usersResponse)

// 		// Assert that the mock service was called as expected
// 		mockService.AssertCalled(t, "GetAllUser", 1, 10, "", 1, 1, 0, "desc")
// 	})
// }

// func (m *MockService) GetUserByID(id uint) (models.User, error) {
// 	args := m.Called(id)
// 	return args.Get(0).(models.User), args.Error(1)
// }

// func TestGetUserByID(t *testing.T) {
// 	// Create a new Fiber instance
// 	app := fiber.New()

// 	// Create a mock service
// 	mockService := new(MockService)

// 	// Set up the route with the handler function
// 	app.Get("/v1/user/detail/:id", func(c *fiber.Ctx) error {
// 		id, _ := c.ParamsInt("id")
// 		user, err := mockService.GetUserByID(uint(id))
// 		if err != nil {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 				"error": "Data not found",
// 			})
// 		}
// 		return c.JSON(user)
// 	})

// 	t.Run("User Found", func(t *testing.T) {
// 		// Expected user for the test
// 		expectedUser := models.User{
// 			Username: "Admin",
// 			RoleID:   1,
// 			Email:    "admin@example.com",
// 		}

// 		// Set up the mock to return the expected user when GetUserByID is called
// 		mockService.On("GetUserByID", uint(1)).Return(expectedUser, nil)

// 		// Create a new HTTP request
// 		req := httptest.NewRequest(http.MethodGet, "/v1/user/detail/1", nil)

// 		// Send the request to the app
// 		resp, err := app.Test(req, -1)
// 		assert.NoError(t, err)

// 		// Check the status code
// 		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

// 		// Check the response body
// 		var userResponse models.User
// 		err = json.NewDecoder(resp.Body).Decode(&userResponse)
// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedUser, userResponse)

// 		// Assert that the mock service was called as expected
// 		mockService.AssertCalled(t, "GetUserByID", uint(1))
// 	})

// 	t.Run("User Not Found", func(t *testing.T) {
// 		// Set up the mock to return an error for a non-existent user
// 		mockService.On("GetUserByID", uint(2)).Return(models.User{}, fmt.Errorf("Data not found"))

// 		// Create a new HTTP request
// 		req := httptest.NewRequest(http.MethodGet, "/v1/user/detail/2", nil)

// 		// Send the request to the app
// 		resp, err := app.Test(req, -1)
// 		assert.NoError(t, err)

// 		// Check the status code
// 		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

// 		// Check the response body
// 		var errorResponse map[string]string
// 		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "Data not found", errorResponse["error"])

// 		// Assert that the mock service was called as expected
// 		mockService.AssertCalled(t, "GetUserByID", uint(2))
// 	})
// }
