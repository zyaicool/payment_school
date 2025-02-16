package controllers

import (
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SchoolYearServiceInterface interface {
	GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, userId int) (response.SchoolYearListResponse, error)
	CreateSchoolYear(schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error)
	UpdateSchoolYear(id uint, schoolYearRequest *request.SchoolYearCreateUpdateRequest, userID int) (*models.SchoolYear, error)
	DeleteSchoolYear(id uint, userID int) (*models.SchoolYear, error)
	GetSchoolYearByID(id uint) (*response.SchoolYearDetailResponse, error)
}

type SchoolYearController struct {
	schoolYearService services.SchoolYearServiceInterface
}

func NewSchoolYearController(schoolYearService services.SchoolYearServiceInterface) *SchoolYearController {
	return &SchoolYearController{schoolYearService: schoolYearService}
}

// @Summary Get All School Years
// @Description Retrieves a paginated list of school years with optional search, sorting, and filters
// @Tags SchoolYear
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of records per page" default(0)
// @Param search query string false "Search term for school year"
// @Param sortBy query string false "Sort by field" default("")
// @Param sortOrder query string false "Sort order" default("asc")
// @Success 200 {array} response.SchoolYearListResponse "List of school years"
// @Failure 500 {object} map[string]interface{} "Failed to fetch data"
// @Router /api/v1/schoolYear/getAllSchoolYear [get]
func (schoolYearController *SchoolYearController) GetAllSchoolYear(c *fiber.Ctx) error {

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	listSchoolYear, err := schoolYearController.schoolYearService.GetAllSchoolYear(page, limit, search, sortBy, sortOrder, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data: " + err.Error(),
		})
	}
	return c.JSON(listSchoolYear)
}

// @Summary Get School Year by ID
// @Description Retrieves a school year by its ID
// @Tags SchoolYear
// @Accept json
// @Produce json
// @Param id path int true "School Year ID"
// @Success 200 {object} response.SchoolYearDetailResponse "School year data"
// @Failure 404 {object} map[string]interface{} "School year not found"
// @Router /api/v1/schoolYear/detail/{id} [get]
func (s *SchoolYearController) GetDataSchoolYear(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	schoolYear, err := s.schoolYearService.GetSchoolYearByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(schoolYear)
}

// @Summary Create a New School Year
// @Description Creates a new school year with the provided details
// @Tags SchoolYear
// @Accept json
// @Produce json
// @Param body body request.SchoolYearCreateUpdateRequest true "School year details"
// @Success 200 {object} models.SchoolYear "Message and created school year data"
// @Failure 400 {object} map[string]interface{} "Bad request - Missing required fields"
// @Failure 500 {object} map[string]interface{} "Failed to create school year"
// @Router /api/v1/schoolYear/create [post]
func (schoolYearController *SchoolYearController) CreateSchoolYear(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolYearRequest *request.SchoolYearCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&schoolYearRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createSchoolYear, err := schoolYearController.schoolYearService.CreateSchoolYear(schoolYearRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createSchoolYear,
	})
}

// @Summary Update a School Year
// @Description Updates an existing school year with the provided details
// @Tags SchoolYear
// @Accept json
// @Produce json
// @Param id path int true "School Year ID"
// @Param body body request.SchoolYearCreateUpdateRequest true "Updated school year details"
// @Success 200 {object} models.SchoolYear "Message and updated school year data"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid input"
// @Failure 500 {object} map[string]interface{} "Failed to update school year"
// @Router /api/v1/schoolYear/update/{id} [put]
func (schoolYearController *SchoolYearController) UpdateSchoolYear(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolYearRequest *request.SchoolYearCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&schoolYearRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedSchoolYear, err := schoolYearController.schoolYearService.UpdateSchoolYear(uint(id), schoolYearRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedSchoolYear,
	})
}

// @Summary Delete a School Year
// @Description Deletes a school year by its ID
// @Tags SchoolYear
// @Accept json
// @Produce json
// @Param id path int true "School Year ID"
// @Success 200 {object} models.SchoolYear "Message indicating successful deletion"
// @Failure 500 {object} map[string]interface{} "Failed to delete school year"
// @Router /api/v1/schoolYear/delete/{id} [delete]
func (schoolYearController *SchoolYearController) DeleteSchoolYear(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = schoolYearController.schoolYearService.DeleteSchoolYear(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}
