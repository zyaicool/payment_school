package controllers

import (
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SchoolGradeController struct {
	schoolGradeService services.SchoolGradeService
}

func NewSchoolGradeController(schoolGradeService services.SchoolGradeService) *SchoolGradeController {
	return &SchoolGradeController{schoolGradeService: schoolGradeService}
}

// @Summary Get All School Grades
// @Description Retrieves a list of school grades with pagination and search options
// @Tags SchoolGrade
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of records per page" default(0)
// @Param search query string false "Search term"
// @Success 200 {array} response.SchoolGradeListResponse "List of school grades"
// @Failure 500 {object} map[string]interface{} "Failed to fetch data"
// @Router /api/v1/schoolGrade/getAllSchoolGrade [get]
func GetAllSchoolGrade(c *fiber.Ctx) error {

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")

	listSchoolGrade, err := services.GetAllSchoolGrade(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(listSchoolGrade)
}

// @Summary Get School Grade by ID
// @Description Retrieves a specific school grade by ID
// @Tags SchoolGrade
// @Accept json
// @Produce json
// @Param id path int true "School Grade ID"
// @Success 200 {object} models.SchoolGrade "School grade data"
// @Failure 404 {object} map[string]interface{} constants.DataNotFoundMessage
// @Router /api/v1/schoolGrade/detail/{id} [get]
func (schoolGradeController *SchoolGradeController) GetDataSchoolGrade(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	schoolGrade, err := schoolGradeController.schoolGradeService.GetSchoolGradeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(schoolGrade)
}

// @Summary Create a New School Grade
// @Description Creates a new school grade with provided data
// @Tags SchoolGrade
// @Accept json
// @Produce json
// @Param body body request.SchoolGradeCreateUpdateRequest true "School grade details"
// @Success 200 {object} models.SchoolGrade "Message and created school grade data"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to create school grade"
// @Router /api/v1/schoolGrade/create [post]
func (schoolGradeController *SchoolGradeController) CreateSchoolGrade(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolGradeRequest *request.SchoolGradeCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&schoolGradeRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createSchoolGrade, err := schoolGradeController.schoolGradeService.CreateSchoolGrade(schoolGradeRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createSchoolGrade,
	})
}

// @Summary Update a School Grade
// @Description Updates an existing school grade with provided data
// @Tags SchoolGrade
// @Accept json
// @Produce json
// @Param id path int true "School Grade ID"
// @Param body body request.SchoolGradeCreateUpdateRequest true "Updated school grade details"
// @Success 200 {object} models.SchoolGrade "Message and updated school grade data"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to update school grade"
// @Router /api/v1/schoolGrade/update/{id} [put]
func (schoolGradeController *SchoolGradeController) UpdateSchoolGrade(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolGradeRequest *request.SchoolGradeCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&schoolGradeRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedSchoolGrade, err := schoolGradeController.schoolGradeService.UpdateSchoolGrade(uint(id), schoolGradeRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedSchoolGrade,
	})
}

// @Summary Delete a School Grade
// @Description Deletes a specific school grade by ID
// @Tags SchoolGrade
// @Accept json
// @Produce json
// @Param id path int true "School Grade ID"
// @Success 200 {object} models.SchoolGrade "Message indicating successful deletion"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to delete school grade"
// @Router /api/v1/schoolGrade/delete/{id} [delete]
func (schoolGradeController *SchoolGradeController) DeleteSchoolGrade(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = schoolGradeController.schoolGradeService.DeleteSchoolGrade(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}
