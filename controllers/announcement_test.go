package controllers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/repositories"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAnnouncementService struct {
	mock.Mock
}

func (m *MockAnnouncementService) GetAnnouncementTypes() ([]repositories.AnnouncementType, error) {
	args := m.Called()
	return args.Get(0).([]repositories.AnnouncementType), args.Error(1)
}

func (m *MockAnnouncementService) CreateAnnouncement(announcementRequest *request.AnnouncementCreateRequest, userID int) (*models.Announcements, error) {
	args := m.Called(announcementRequest, userID)
	return args.Get(0).(*models.Announcements), args.Error(1)
}

func (m *MockAnnouncementService) DeleteAnnouncement(id uint, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockAnnouncementService) UpdateAnnouncement(announcementID int, updateRequest *request.AnnouncementUpdateRequest, userID int) (*models.Announcements, error) {
	args := m.Called(announcementID, updateRequest, userID)
	return args.Get(0).(*models.Announcements), args.Error(1)
}

func (m *MockAnnouncementService) GetListAnnouncement(page, limit int, search, sortBy, sortOrder, announcementType string, userID int) (response.GetAnnouncementListResponse, error) {
	args := m.Called(page, limit, search, sortBy, sortOrder, announcementType, userID)
	return args.Get(0).(response.GetAnnouncementListResponse), args.Error(1)
}

func (m *MockAnnouncementService) GetAnnouncementByID(id int) (response.AnnouncementResponse, error) {
	args := m.Called(id)
	return args.Get(0).(response.AnnouncementResponse), args.Error(1)
}

func TestCreateAnnouncement_NormalCase(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
		return c.Next()
	})

	mockService := new(MockAnnouncementService)
	controller := NewAnnouncementController(mockService)

	mockService.On("CreateAnnouncement", mock.AnythingOfType("*request.AnnouncementCreateRequest"), 1).
		Return(&models.Announcements{
			Title: "tes",
		}, nil)

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("schoolId", "1")
	writer.WriteField("heroImage", "")
	writer.WriteField("title", "tes")
	writer.WriteField("description", "create announcement")
	writer.WriteField("type", "AT02")
	writer.WriteField("eventDate", "")
	writer.Close()

	app.Post("/create", controller.CreateAnnouncement)
	req := httptest.NewRequest("POST", "/create", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the response message and data
	assert.Equal(t, "Announcement created successfully.", responseBody["message"])
	data := responseBody["data"].(map[string]interface{})
	assert.Equal(t, "tes", data["title"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

func TestCreateAnnouncement_schoolId_invalid(t *testing.T) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
		return c.Next()
	})

	mockService := new(MockAnnouncementService)
	controller := NewAnnouncementController(mockService)

	// Create form-data request body with invalid schoolId (non-numeric string)
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("schoolId", "invalidID") // Non-numeric to trigger strconv.Atoi error
	writer.WriteField("heroImage", "")
	writer.WriteField("title", "tes")
	writer.WriteField("description", "create announcement")
	writer.WriteField("type", "AT02")
	writer.WriteField("eventDate", "")
	writer.Close()

	app.Post("/create", controller.CreateAnnouncement)
	req := httptest.NewRequest("POST", "/create", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var responseBody fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Assert the response message and data
	assert.Equal(t, "Invalid school ID", responseBody["error"])

	// Assert that the mock expectations were met
	mockService.AssertExpectations(t)
}

// func TestCreateAnnouncement_PositiveCase_FileUpload(t *testing.T) {
// 	app := fiber.New()

// 	app.Use(func(c *fiber.Ctx) error {
// 		c.Locals("user", jwt.MapClaims{"user_id": float64(1)})
// 		return c.Next()
// 	})

// 	mockService := new(MockAnnouncementService)
// 	controller := NewAnnouncementController(mockService)

// 	// Mock the service to return a successful response
// 	mockService.On("CreateAnnouncement", mock.AnythingOfType("*request.AnnouncementCreateRequest"), 1).
// 		Return(&models.Announcements{
// 			Title: "tes",
// 		}, nil)

// 	// Mock the utility function to simulate successful image upload
// 	utilities.UploadImage = func(file *multipart.FileHeader, c *fiber.Ctx, dir, fieldName string) (string, error) {
// 		return "/path/to/uploaded/image.jpg", nil
// 	}

// 	// Create form-data request body with a file
// 	var buffer bytes.Buffer
// 	writer := multipart.NewWriter(&buffer)
// 	writer.WriteField("schoolId", "1")
// 	writer.WriteField("title", "Test Announcement")
// 	writer.WriteField("description", "This is a test announcement.")
// 	writer.WriteField("type", "AT02")
// 	writer.WriteField("eventDate", "2025-01-01")

// 	// Add a file to the form-data
// 	part, _ := writer.CreateFormFile("heroImage", "test.jpg")
// 	part.Write([]byte("test file content"))
// 	writer.Close()

// 	app.Post("/create", controller.CreateAnnouncement)
// 	req := httptest.NewRequest("POST", "/create", &buffer)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	resp, err := app.Test(req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

// 	var responseBody map[string]interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&responseBody)
// 	assert.NoError(t, err)

// 	// Assert the response message and data
// 	assert.Equal(t, "Announcement created successfully.", responseBody["message"])
// 	data := responseBody["data"].(map[string]interface{})
// 	assert.Equal(t, "tes", data["title"])

// 	// Assert that the mock expectations were met
// 	mockService.AssertExpectations(t)
// }
