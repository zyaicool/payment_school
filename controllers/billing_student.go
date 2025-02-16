package controllers

import (
	"fmt"
	"strconv"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/models"
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BillingStudentController struct {
	billingStudentService services.BillingStudentServiceInterface
}

func NewBillingStudentController(billingStudentService services.BillingStudentServiceInterface) *BillingStudentController {
	return &BillingStudentController{billingStudentService: billingStudentService}
}

// @Summary Get All Billing Students
// @Description Retrieve a paginated and filtered list of billing students
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Param search query string false "Search query for billing students"
// @Param studentId query int false "Filter by student ID"
// @Param schoolClassId query int false "Filter by school class ID"
// @Param paymentType query string false "Filter by payment type"
// @Param schoolGradeId query int false "Filter by school grade ID"
// @Param schoolId query int false "Filter by school ID"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc or desc, default: asc)"
// @Success 200 {array} models.BillingStudent "List of billing students"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/billingStudent/getAllBillingStudent [get]
func (billingStudentController BillingStudentController) GetAllBillingStudent(c *fiber.Ctx) error {
	// err := utilities.CheckAccessExceptUserOrtu(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	userID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")
	studentID := c.QueryInt("studentId", 0)
	schoolClassID := c.QueryInt("schoolClassId", 0)

	paymentType := c.Query("paymentType")
	schoolGradeID := c.QueryInt("schoolGradeId", 0)
	schoolID := c.QueryInt("schoolId", 0)
	// sort := c.Query("sort", "desc") // Default to "desc"
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")

	listBilling, err := billingStudentController.billingStudentService.GetAllBillingStudent(page, limit, search, studentID, roleID, schoolGradeID, paymentType, schoolID, userID, sortBy, sortOrder, schoolClassID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(listBilling)
}

func CheckAccessToBillingStudent(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := services.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID == 3 {
		return fmt.Errorf("User can't access this page")
	}

	return nil
}

// @Summary Get Billing Detail by ID
// @Description Retrieve detailed billing information by student ID
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param studentId query int true "Student ID"
// @Success 200 {object} models.BillingDetail "Detailed billing information"
// @Failure 404 {object} map[string]interface{} constants.DataNotFoundMessage
// @Router /api/v1/billingStudent/detailBillingStudent [get]
func (billingStudentController BillingStudentController) GetDetailBillingID(c *fiber.Ctx) error {
	// err := utilities.CheckAccessExceptUserOrtu(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	studentId := c.QueryInt("studentId")

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := uint(userClaims["user_id"].(float64))
	roleId := uint(userClaims["role_id"].(float64))
	schoolId := uint(userClaims["school_id"].(float64))

	var user = models.User{
		Master: models.Master{
			ID: userID,
		},
		RoleID: roleId,
		UserSchool: &models.UserSchool{
			School: &models.School{
				ID: schoolId,
			},
		},
	}

	billing, err := billingStudentController.billingStudentService.GetDetailBillingIDService(studentId, user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(billing)
}

// @Summary Get Detailed Billing Student
// @Description Retrieve detailed billing student data by ID
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param id path int true "Billing Student ID"
// @Success 200 {object} response.BillingStudentDetailResponse "Billing student details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/billingStudent/detailBillingByBillingStudentId/{id} [get]
func (billingStudentController *BillingStudentController) GetDetailBillingStudent(c *fiber.Ctx) error {
	billingStudentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidBIllingStudentIdMessage,
		})
	}

	billingDetail, err := billingStudentController.billingStudentService.GetDetailBillingStudentByID(billingStudentId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve billing student details",
		})
	}

	return c.JSON(billingDetail)
}

// @Summary Update Billing Student
// @Description Update billing student data
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param id path int true "Billing Student ID"
// @Param updateRequest body request.UpdateBillingStudentRequest true "Billing student update request"
// @Success 200 {object} models.BillingStudent "Updated billing student data"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/billingStudent/update/{id} [put]
func (controller *BillingStudentController) UpdateBillingStudent(c *fiber.Ctx) error {
	billingStudentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidBIllingStudentIdMessage,
		})
	}

	var updateRequest request.UpdateBillingStudentRequest
	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	responseData, err := controller.billingStudentService.UpdateBillingStudentService(billingStudentId, updateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update billing student data",
		})
	}

	return c.JSON(responseData)
}

// @Summary Delete Billing Student
// @Description Delete a billing student record
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param id path int true "Billing Student ID"
// @Success 200 {object} map[string]interface{} "Delete confirmation"
// @Failure 400 {object} map[string]interface{} "Invalid ID or deletion error"
// @Router /api/v1/billingStudent/delete/{id} [delete]
func (controller *BillingStudentController) DeleteBillingStudent(c *fiber.Ctx) error {
	billingStudentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.InvalidBIllingStudentIdMessage,
		})
	}
	userClaims := c.Locals("user").(jwt.MapClaims)
	deletedBy := int(userClaims["user_id"].(float64))

	err = controller.billingStudentService.DeleteBillingStudentService(billingStudentId, deletedBy)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to delete billing student(check if billing student exists)",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Billing student deleted successfully",
	})
}

// @Summary Create Billing Student
// @Description Create a new billing student record
// @Tags Billing Students
// @Accept json
// @Produce json
// @Param createRequest body request.CreateBillingStudentRequest true "Billing student creation request"
// @Success 200 {object} map[string]interface{} "Created billing student details"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/billingStudent/create [post]
func (controller *BillingStudentController) CreateBillingStudent(c *fiber.Ctx) error {
	var createRequest request.CreateBillingStudentRequest
	if err := c.BodyParser(&createRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.FailedToParseRequestBodyMessage,
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	billingStudent, err := controller.billingStudentService.CreateBillingStudent(createRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    billingStudent,
		"message": "Billing student created successfully",
	})
}
