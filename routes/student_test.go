package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStudentController is a mock of the StudentController
type MockStudentController struct {
	mock.Mock
}

func (m *MockStudentController) GetAllStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) GetStudentByID(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) CreateStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) UpdateStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) DeleteStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) CreateLinkToUser(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) GetStudentByUserId(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) DownloadFileExcelFormatForStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) UploadStudents(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) GetStudentImage(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) UploadImageStudent(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockStudentController) ExportStudentToExcel(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func TestSetupStudentRoutes(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")

	SetupStudentRoutes(api)

	// Test route registration
	stack := app.Stack()
	assert.NotEmpty(t, stack)

	// Verify specific routes
	expectedRoutes := []string{
		"GET /api/student/getAllStudent",
		"GET /api/student/detail/:id",
		"POST /api/student/create",
		"PUT /api/student/update/:id",
		"DELETE /api/student/delete/:id",
		"POST /api/student/linkToUser",
		"GET /api/student/user/:user_id",
		"GET /api/student/getFileExcel",
		"POST /api/student/upload",
		"GET /api/student/image/:id",
		"POST /api/student/image/:id",
		"GET /api/student/exportExcel",
	}

	for _, route := range expectedRoutes {
		found := false
		for _, routeStack := range stack {
			for _, r := range routeStack {
				if r.Method+" "+r.Path == route {
					found = true
					break
				}
			}
		}
		assert.True(t, found, "Route %s should be registered", route)
	}
}
