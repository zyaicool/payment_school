package controllers

import (
	"schoolPayment/constants"
	"schoolPayment/dtos/request"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SchoolMajorControllerInterface interface {
	GetAllSchoolMajor(c *fiber.Ctx) error
	CreateSchoolMajor(c *fiber.Ctx) error
}

type SchoolMajorController struct {
	schoolMajorService services.SchoolMajorService
}

func NewSchoolMajorController(schoolMajorService services.SchoolMajorService) SchoolMajorControllerInterface {
	return &SchoolMajorController{schoolMajorService: schoolMajorService}
}

// @Summary Get All School Majors
// @Description Retrieves a list of school majors with optional search filter
// @Tags SchoolMajor
// @Accept json
// @Produce json
// @Param search query string false "Search term for school major"
// @Success 200 {array} response.SchoolMajorResponse "List of school majors"
// @Failure 500 {object} map[string]interface{} "Failed to fetch data"
// @Router /api/v1/schoolMajor/getListSchoolMajor [get]
func (schoolMajorController *SchoolMajorController) GetAllSchoolMajor(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64)) 

	search := c.Query("search")

	major, err := schoolMajorController.schoolMajorService.GetAllSchoolMajorService(search, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(major)
}

// @Summary Create a New School Major
// @Description Creates a new school major with the provided details
// @Tags SchoolMajor
// @Accept json
// @Produce json
// @Param body body request.SchoolMajorCreate true "School major details"
// @Success 200 {object} models.SchoolMajor "Message and created school major data"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing required fields"
// @Failure 500 {object} map[string]interface{} "Failed to create school major"
// @Router /api/v1/schoolMajor/create [post]
func (schoolMajorController *SchoolMajorController) CreateSchoolMajor(c *fiber.Ctx) error {
	var major *request.SchoolMajorCreate
	var userID int = 0

	// Get user information from the token claims stored in context
	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		// If the user claims exist and are of type jwt.MapClaims, extract the user ID.
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID) 
		}
	}

	if err := c.BodyParser(&major); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	if major.SchoolMajorName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SchoolMajorName is required",
		})
	}
	if major.SchoolID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SchoolID is required",
		})
	}

	createMajor, err := schoolMajorController.schoolMajorService.CreateSchoolMajorService(major, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Data berhasil disimpan.",
		"data":    createMajor,
	})
}
