package controllers

import (
	"strconv"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New()

type SchoolController struct {
	schoolService services.SchoolServiceInterface
}

func NewSchoolController(schoolService services.SchoolServiceInterface) *SchoolController {
	return &SchoolController{schoolService: schoolService}
}

// @Summary Get All Schools
// @Description Get a list of all schools with pagination, search, and sorting options
// @Tags Schools
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param search query string false "Search term"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order" enum(asc,desc) default(asc)
// @Success 200 {array} response.SchoolListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/getAllSchool [get]
func (schoolController *SchoolController) GetAllSchool(c *fiber.Ctx) error {

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")

	listSchool, err := schoolController.schoolService.GetAllSchoolList(page, limit, search, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(listSchool)
}

// @Summary Get School Data by ID
// @Description Get details of a school by its ID
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Success 200 {object} response.DetailSchoolResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/detail/{id} [get]
func (schoolController *SchoolController) GetDataSchool(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")
	school, err := schoolController.schoolService.GetSchoolByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(school)
}

// @Summary Create a New School
// @Description Create a new school with all required fields and optional logo and letterhead files
// @Tags Schools
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param npsn formData int true "NPSN"
// @Param schoolName formData string true "School Name"
// @Param schoolProvince formData string true "School Province"
// @Param schoolCity formData string true "School City"
// @Param schoolPhone formData string true "School Phone"
// @Param schoolAddress formData string true "School Address"
// @Param schoolMail formData string true "School Mail"
// @Param schoolFax formData string true "School Fax"
// @Param schoolGradeId formData int true "School Grade ID"
// @Param schoolLogo formData file false "School Logo"
// @Param schoolLetterhead formData file false "School Letterhead"
// @Success 200 {object} models.School
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/create [post]
func (schoolController *SchoolController) CreateSchool(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	npsn, err := strconv.Atoi(c.FormValue("npsn"))
	if err != nil {
		npsn = 0 // Set default value if parsing fails
	}

	schoolGradeId, err := strconv.Atoi(c.FormValue("schoolGradeId"))
	if err != nil {
		schoolGradeId = 0 // Set default value if parsing fails
	}

	schoolRequest := new(request.SchoolCreateUpdateRequest)
	schoolRequest.Npsn = npsn
	schoolRequest.SchoolName = c.FormValue("schoolName")
	schoolRequest.SchoolProvince = c.FormValue("schoolProvince")
	schoolRequest.SchoolCity = c.FormValue("schoolCity")
	schoolRequest.SchoolPhone = c.FormValue("schoolPhone")
	schoolRequest.SchoolAddress = c.FormValue("schoolAddress")
	schoolRequest.SchoolMail = c.FormValue("schoolMail")
	schoolRequest.SchoolFax = c.FormValue("schoolFax")
	schoolRequest.SchoolGradeID = schoolGradeId

	fileLogo, err := c.FormFile("schoolLogo")
	if err == nil { // Process only if a file is provided
		newFileLogoName, err := utilities.UploadImage(fileLogo, c, "school_logo", "schoolLogo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		schoolRequest.SchoolLogo = newFileLogoName
	}

	fileLetterhead, err := c.FormFile("schoolLetterhead")
	if err == nil { // Process only if a file is provided
		newFileLogoName, err := utilities.UploadImage(fileLetterhead, c, "school_letterhead", "schoolLetterhead")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		schoolRequest.SchoolLetterhead = newFileLogoName
	}

	if err := validate.Struct(schoolRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	createSchool, err := schoolController.schoolService.CreateSchool(schoolRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createSchool,
	})
}

// @Summary Update School Data by ID
// @Description Update an existing school's information by its ID
// @Tags Schools
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "School ID"
// @Param npsn formData int false "NPSN"
// @Param schoolName formData string false "School Name"
// @Param schoolProvince formData string false "School Province"
// @Param schoolCity formData string false "School City"
// @Param schoolPhone formData string false "School Phone"
// @Param schoolAddress formData string false "School Address"
// @Param schoolMail formData string false "School Mail"
// @Param schoolFax formData string false "School Fax"
// @Param schoolGradeId formData int false "School Grade ID"
// @Param schoolLogo formData file false "School Logo"
// @Param schoolLetterhead formData file false "School Letterhead"
// @Success 200 {object} models.School
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/update/{id} [put]
func (schoolController *SchoolController) UpdateSchool(c *fiber.Ctx) error {
	// Check access to the school
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Retrieve user ID and school ID from claims and params
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse form data into the DTO
	npsn, err := strconv.Atoi(c.FormValue("npsn"))
	if err != nil {
		npsn = 0 // Set default value if parsing fails
	}

	schoolGradeId, err := strconv.Atoi(c.FormValue("schoolGradeId"))
	if err != nil {
		schoolGradeId = 0 // Set default value if parsing fails
	}

	// Populate the SchoolCreateUpdateRequest struct
	schoolRequest := new(request.SchoolCreateUpdateRequest)
	schoolRequest.Npsn = npsn
	schoolRequest.SchoolName = c.FormValue("schoolName")
	schoolRequest.SchoolProvince = c.FormValue("schoolProvince")
	schoolRequest.SchoolCity = c.FormValue("schoolCity")
	schoolRequest.SchoolPhone = c.FormValue("schoolPhone")
	schoolRequest.SchoolAddress = c.FormValue("schoolAddress")
	schoolRequest.SchoolMail = c.FormValue("schoolMail")
	schoolRequest.SchoolFax = c.FormValue("schoolFax")
	schoolRequest.SchoolGradeID = schoolGradeId

	// Handle file upload for school logo
	fileLogo, err := c.FormFile("schoolLogo")
	if err == nil { // Process only if a file is provided
		newFileLogoName, err := utilities.UploadImage(fileLogo, c, "school_logo", "schoolLogo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		schoolRequest.SchoolLogo = newFileLogoName
	}

	fileLetterhead, err := c.FormFile("schoolLetterhead")
	if err == nil { // Process only if a file is provided
		newFileLogoName, err := utilities.UploadImage(fileLetterhead, c, "school_letterhead", "schoolLetterhead")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		schoolRequest.SchoolLetterhead = newFileLogoName
	}

	// Call the service to update the school data
	updatedSchool, err := schoolController.schoolService.UpdateSchool(uint(id), schoolRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil diperbarui.",
		"data":    updatedSchool,
	})
}

// @Summary Delete School by ID
// @Description Delete a school by its ID
// @Tags Schools
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "School ID"
// @Success 200 {object} models.School
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/delete/{id} [delete]
func (schoolController *SchoolController) DeleteSchool(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = schoolController.schoolService.DeleteSchool(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}

// @Summary Get All Onboarding Schools
// @Description Get a list of all onboarding schools with a search option
// @Tags Schools
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Success 200 {array} []response.SchoolOnBoardingResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/school/getAllOnboarding [get]
func (schoolController *SchoolController) GetAllOnboardingSchools(ctx *fiber.Ctx) error {

	search := ctx.Query("search", "")

	schoolsOnboarding, err := schoolController.schoolService.GetAllOnboardingSchools(search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch schools",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schoolsOnboarding,
	})
}
