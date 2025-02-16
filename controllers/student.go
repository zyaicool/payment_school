package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"schoolPayment/configs"
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"gorm.io/gorm"
)

var (
	studentControllerInstance *StudentController
	studentControllerOnce     sync.Once
)

type StudentController struct {
	studentService services.StudentService
	DB             *gorm.DB
}

func GetStudentController() *StudentController {
	studentControllerOnce.Do(func() {
		studentService := services.GetStudentService()
		studentControllerInstance = &StudentController{
			studentService: studentService,
			DB:             configs.DB,
		}
	})
	return studentControllerInstance
}

// @Summary Get All Students
// @Description Retrieve a list of all students with various filters
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of records per page" default(10)
// @Param search query string false "Search term"
// @Param searchNis query string false "Search by NIS"
// @Param status query string false "Student status"
// @Param gradeId query int false "Grade ID"
// @Param yearId query int false "Year ID"
// @Param schoolId query int false "School ID"
// @Param classId query int false "Class ID"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order" enum([asc, desc])
// @Param studentId query int false "Student ID"
// @Success 200 {array} response.StudentListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/getAllStudent [get]
func (studentController *StudentController) GetAllStudent(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")
	searchNis := c.Query("searchNis", "")
	status := c.Query("status")
	gradeID := c.QueryInt("gradeId", 0)
	yearId := c.QueryInt("yearId", 0)
	schoolID := c.QueryInt("schoolId", 0)
	classID := c.QueryInt("classId", 0)
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")
	studentId := c.QueryInt("studentId", 0)

	var isActive *bool
	isActiveStr := c.Query("isActive", "")
	if isActiveStr != "" {
		val := strings.ToLower(isActiveStr) == "true"
		isActive = &val
	}

	students, err := studentController.studentService.GetAllStudent(page, limit, search, userID, status, gradeID, yearId, schoolID, searchNis, classID, sortBy, sortOrder, studentId, isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(students)
}

// @Summary Get Student By ID
// @Description Retrieve student details by their ID
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Student ID"
// @Success 200 {object} response.DetailStudentResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/detail/{id} [get]
func (studentController *StudentController) GetStudentByID(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	id, _ := c.ParamsInt("id")
	student, err := studentController.studentService.GetStudentByID(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(student)
}

// @Summary Create New Student
// @Description Create a new student record
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param body body request.CreateStudentRequest true "Student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/create [post]
func (studentController *StudentController) CreateStudent(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var createStudentRequest request.CreateStudentRequest
	var userId int = 0

	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userId = int(userClaimID)
		}
	}

	// Parse JSON into the DTO
	if err := c.BodyParser(&createStudentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fmt.Sprintf("Failed to parse request body: %v", err),
			"message": "Please ensure the request body is valid JSON and matches the expected format",
		})
	}

	createStudent, err := studentController.studentService.CreateStudent(createStudentRequest, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createStudent,
	})
}

// @Summary Update Student Data
// @Description Update an existing student's data
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Student ID"
// @Param body body request.UpdateStudentRequest true "Updated student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/update/{id} [put]
func (studentController *StudentController) UpdateStudent(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil user ID dari token JWT
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Ambil student ID dari parameter URL
	id, _ := c.ParamsInt("id")
	var updateStudentRequest request.UpdateStudentRequest

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&updateStudentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   fmt.Sprintf("Failed to parse request body: %v", err),
			"message": "Please ensure the request body is valid JSON and matches the expected format",
		})
	}

	updatedStudentPointer, err := studentController.studentService.UpdateStudent(uint(id), updateStudentRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedStudentPointer,
	})
}

// @Summary Delete Student
// @Description Delete a student record by ID
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/delete/{id} [delete]
func (studentController *StudentController) DeleteStudent(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = studentController.studentService.DeleteStudentService(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}

func GetUserStudent(c *fiber.Ctx) error {
	panic("implement me")
}

func (studentController *StudentController) CreateLinkToUser(c *fiber.Ctx) error {
	id_user := c.QueryInt("user_id")
	studentID := c.QueryInt("student_id")

	userClaims := c.Locals("user").(jwt.MapClaims)
	auth := int(userClaims["user_id"].(float64))

	createUserStudent, err := studentController.studentService.CreateUserStudentService(id_user, studentID, auth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(createUserStudent)
}

// @Summary Get Student By User ID
// @Description Retrieve a student by user ID
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param user_id path int true "User ID"
// @Success 200 {object} []models.Student
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/user/{user_id} [get]
func (studentController *StudentController) GetStudentByUserId(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("user_id")
	student, err := studentController.studentService.GetStudentByUserIDService(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(student)
}

func (studentController *StudentController) GetStudentImage(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	student, err := studentController.studentService.GetStudentByID(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}

	// Check if the student has an image
	if student.Image == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Image not found for student",
		})
	}

	// Generate the URL for the image
	// baseUrl := "http://127.0.0.1:8081"
	baseUrl := "http://localhost:8081"
	imageUrl := fmt.Sprintf("%s/upload/%s", baseUrl, student.Image)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"imageUrl": imageUrl,
	})
}

// @Summary Download Excel Format for Students
// @Description Download an Excel file template for students
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Success 200 {file} file "Excel file template"
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/getFileExcel [get]
func (studentController *StudentController) DownloadFileExcelFormatForStudent(c *fiber.Ctx) error {
	// Check access for users except role 'Ortu'
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate the Excel file
	buffer, err := services.GenerateFileExcelForStudent(c, studentController.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStream(buffer)
}

// @Summary Upload Students
// @Description Upload student data from a CSV or Excel file
// @Tags Students
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param file formData file true "Student data file"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/upload [post]
func (studentController *StudentController) UploadStudents(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	// Check file type (CSV or Excel)
	extension := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	// fmt.Println(extension, file)
	switch extension {
	case ".xlsx", ".xls":
		return studentController.studentService.HandleExcelFileStudent(file, c, userID)
	case ".csv":
		return studentController.studentService.HandleCSVFileStudent(file, c, userID)
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unsupported file type",
		})
	}
}

func (studentController *StudentController) UploadImageStudent(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	var request request.CreateStudentImageRequest
	var userId int = 0
	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userId = int(userClaimID)
		}
	}

	// ini logic untuk upload image, nanti di pakai ketika di butuhkan
	file, err := c.FormFile("image")
	if err != nil {
		// Log or return the error if image is missing or not uploaded correctly
		fmt.Println("Error retrieving file:", err)
	} else {
		epoch := time.Now().Unix()

		extension := filepath.Ext(file.Filename)

		newFileName := fmt.Sprintf("%d%s", epoch, extension)

		// Ensure the upload directory exists
		uploadDir := "./upload/"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, 0o755) // Create the directory with necessary permissions
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Unable to create upload directory",
				})
			}
		}

		filePath := filepath.Join(uploadDir, newFileName)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to save image",
			})
		}

		request.Image = newFileName
	}

	createStudent, err := studentController.studentService.CreateImageStudentService(request, uint(id), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createStudent,
	})
}

// @Summary Export Students to Excel
// @Description Export student data to an Excel file
// @Tags Students
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param search query string false "Search term"
// @Param searchNis query string false "Search by NIS"
// @Param gradeId query int false "Grade ID"
// @Param yearId query int false "Year ID"
// @Param status query string false "Student status"
// @Param schoolId query int false "School ID"
// @Success 200 {file} file "Exported Excel file"
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/student/exportExcel [get]
func (studentController *StudentController) ExportStudentToExcel(c *fiber.Ctx) error {
	// Get query parameters
	search := c.Query("search", "")
	searchNis := c.Query("searchNis", "")
	gradeId, _ := strconv.Atoi(c.Query("gradeId", "0"))
	yearId, _ := strconv.Atoi(c.Query("yearId", "0"))
	status := c.Query("status", "")
	schoolId, _ := strconv.Atoi(c.Query("schoolId", "0"))

	var isActive *bool
	isActiveStr := c.Query("isActive", "")
	if isActiveStr != "" {
		val := strings.ToLower(isActiveStr) == "true"
		isActive = &val
	}

	// Get user ID from JWT token
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Call service to generate Excel file
	exportData, err := studentController.studentService.ExportStudentToExcel(search, searchNis, gradeId, yearId, status, schoolId, userID, isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Format waktu dan nama file
	currentTime := time.Now()
	formattedTime := currentTime.Format("02-01-2006 15.04")
	filename := fmt.Sprintf("Data Siswa %s %s.xlsx", strings.ReplaceAll(exportData.SchoolName, " ", "_"), formattedTime)

	// Set header untuk unduhan file
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.SendStream(exportData.Buffer)
}
