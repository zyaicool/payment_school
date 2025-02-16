package controllers

import (
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/models"
	"schoolPayment/repositories"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type SchoolClassController struct {
	schoolClassService services.SchoolClassService
}

func NewSchoolClassController() *SchoolClassController {
	return &SchoolClassController{
		schoolClassService: services.NewSchoolClassService(
			&repositories.SchoolClassRepository{},
			repositories.NewUserRepository(),
		),
	}
}

// @Summary Get All School Classes
// @Description Retrieves a list of school classes with pagination, sorting, and search options
// @Tags SchoolClass
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of records per page" default(0)
// @Param search query string false "Search term"
// @Param sortBy query string false "Sort field"
// @Param sortOrder query string false "Sort order" default(asc)
// @Param showDeletedData query bool false "Show deleted data" default(false)
// @Success 200 {array} response.SchoolClassListResponse "List of school classes"
// @Failure 500 {object} map[string]interface{} "Failed to fetch data"
// @Router /api/v1/schoolClass/getAllSchoolClass [get]
func (schoolClassController *SchoolClassController) GetAllSchoolClass(c *fiber.Ctx) error {
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

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")
	showDeletedData := c.QueryBool("showDeletedData", false)

	listSchoolClass, err := schoolClassController.schoolClassService.GetAllSchoolClass(page, limit, search, sortBy, sortOrder, showDeletedData, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(listSchoolClass)
}

// @Summary Get School Class by ID
// @Description Retrieves a specific school class by ID
// @Tags SchoolClass
// @Accept json
// @Produce json
// @Param id path int true "School Class ID"
// @Success 200 {object} response.SchoolClassDetailResponse "School class data"
// @Failure 404 {object} map[string]interface{} constants.DataNotFoundMessage
// @Router /api/v1/schoolClass/detail/{id} [get]
func (schoolClassController *SchoolClassController) GetDataSchoolClass(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")
	schoolClass, err := schoolClassController.schoolClassService.GetSchoolClassByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(schoolClass)
}

// @Summary Create a New School Class
// @Description Creates a new school class with provided data
// @Tags SchoolClass
// @Accept json
// @Produce json
// @Param body body request.SchoolClassCreateUpdateRequest true "School class details"
// @Success 200 {object} models.SchoolClass "Message and created school class data"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to create school class"
// @Router /api/v1/schoolClass/create [post]
func (schoolClassController *SchoolClassController) CreateSchoolClass(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolClassRequest *request.SchoolClassCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&schoolClassRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createSchoolClass, err := schoolClassController.schoolClassService.CreateSchoolClass(schoolClassRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createSchoolClass,
	})
}

// @Summary Update a School Class
// @Description Updates an existing school class with provided data
// @Tags SchoolClass
// @Accept json
// @Produce json
// @Param id path int true "School Class ID"
// @Param body body request.SchoolClassCreateUpdateRequest true "Updated school class details"
// @Success 200 {object} models.SchoolClass "Message and updated school class data"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to update school class"
// @Router /api/v1/schoolClass/update/{id} [put]
func (schoolClassController *SchoolClassController) UpdateSchoolClass(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var schoolClassRequest *request.SchoolClassCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&schoolClassRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedSchoolClass, err := schoolClassController.schoolClassService.UpdateSchoolClass(uint(id), schoolClassRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedSchoolClass,
	})
}

// @Summary Delete a School Class
// @Description Deletes a specific school class by ID
// @Tags SchoolClass
// @Accept json
// @Produce json
// @Param id path int true "School Class ID"
// @Success 200 {object} models.SchoolClass "Message indicating successful deletion"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to delete school class"
// @Router /api/v1/schoolClass/delete/{id} [delete]
func (schoolClassController *SchoolClassController) DeleteSchoolClass(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = schoolClassController.schoolClassService.DeleteSchoolClass(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
	})
}
