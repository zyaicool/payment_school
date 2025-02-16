package controllers

import (
	"net/http/httptest"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBillingStudentService struct {
	mock.Mock
}

func (m *MockBillingStudentService) GetAllBillingStudent(page int, limit int, search string, studentID int, roleID int, schoolGradeID int, paymentType string, schoolID int, userID int, sortBy string, sortOrder string, schoolClassID int) (response.ListBillingPerStudent, error) {
	args := m.Called(page, limit, search, studentID, roleID, schoolGradeID, paymentType, schoolID, userID, sortBy, sortOrder, schoolClassID)
	return args.Get(0).(response.ListBillingPerStudent), args.Error(1)
}

func (m *MockBillingStudentService) GetDetailBillingIDService(studentId int, user models.User) (*response.BillingStudentByStudentIDResponse, error) {
	args := m.Called(studentId, user)
	return args.Get(0).(*response.BillingStudentByStudentIDResponse), args.Error(1)
}

func (m *MockBillingStudentService) GetDetailBillingStudentByID(billingStudentId int) (*response.BillingStudentDetailResponse, error) {
	args := m.Called(billingStudentId)
	return args.Get(0).(*response.BillingStudentDetailResponse), args.Error(1)
}

func (m *MockBillingStudentService) UpdateBillingStudentService(billingStudentId int, updateRequest request.UpdateBillingStudentRequest) (response.BillingStudentDetailResponse, error) {
	args := m.Called(billingStudentId, updateRequest)
	return args.Get(0).(response.BillingStudentDetailResponse), args.Error(1)
}

func (m *MockBillingStudentService) DeleteBillingStudentService(billingStudentId, deletedBy int) error {
	args := m.Called(billingStudentId, deletedBy)
	return args.Error(0)
}

func (m *MockBillingStudentService) CreateBillingStudent(request request.CreateBillingStudentRequest, userID int) ([]models.BillingStudent, error) {
	args := m.Called(request, userID)
	return args.Get(0).([]models.BillingStudent), args.Error(1)
}


func TestGetBillingStudent_Success(t *testing.T) {
    // Initialize the app
    app := fiber.New()

    // Mock service
    mockService := new(MockBillingStudentService)
    controller := NewBillingStudentController(mockService)

    // Middleware to set mock JWT claims
    app.Use(func(c *fiber.Ctx) error {
        c.Locals("user", jwt.MapClaims{
            "user_id":   float64(1), // Mocked user ID
            "role_id":   float64(4), // Mocked role ID (Kasir, which should pass)
            "school_id": float64(1),
        })
        return c.Next()
    })

    // Mock the GetDetailBillingIDService method to return a valid billing student response
    studentId := 1

    mockBillingResponse := &response.BillingStudentByStudentIDResponse{
        ID:               uint(studentId),
        Nis:              "12345",
        StudentName:      "John Doe",
        SchoolClassName:  "Class A",
        SchoolGradeName:  "Grade 10",
        LatestSchoolYear: "2023/2024",
        ListBilling:      []response.BillingStudentByStudentIDDetailResponse{},
        ListDonation:     []response.DonationBillingResponse{},
    }

    // Mock service behavior using mock.Anything to match any argument
    mockService.On("GetDetailBillingIDService", mock.Anything, mock.Anything).
        Return(mockBillingResponse, nil).Once()

    app.Get("/billing/student/:id", controller.GetDetailBillingID)

    req := httptest.NewRequest("GET", "/billing/student/1", nil)
    req.Header.Set("Authorization", "Bearer mock-token")

    // Send the request and get the response
    resp, err := app.Test(req)

    // Assert no error occurred
    assert.NoError(t, err)

    // Assert the response status is OK (200)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    // Verify that the mock method was called as expected
    mockService.AssertExpectations(t)
}
