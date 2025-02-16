package controllers

import (
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/services"
	"schoolPayment/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AnnouncementController struct {
	announcementService services.AnnouncementService
}

func NewAnnouncementController(announcementService services.AnnouncementService) *AnnouncementController {
	return &AnnouncementController{announcementService: announcementService}
}

// @Summary Create a New Announcement
// @Description Create a new announcement with details and optional hero image
// @Tags Announcements
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param schoolId formData int true "School ID"
// @Param title formData string true "Title of the announcement"
// @Param description formData string false "Description of the announcement"
// @Param type formData string true "Type of the announcement (e.g., Event, General)"
// @Param eventDate formData string false "Event Date in YYYY-MM-DD format"
// @Param heroImage formData file false "Hero image for the announcement"
// @Success 200 {object} map[string]interface{} "Announcement created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/announcement/create [post]
func (announcementController *AnnouncementController) CreateAnnouncement(c *fiber.Ctx) error {
	var userID int = 0

	// Extract user ID from token claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	// Parse form fields
	schoolIDStr := c.FormValue("schoolId")
	title := c.FormValue("title")
	description := c.FormValue("description")
	announcementType := c.FormValue("type")
	eventDateStr := c.FormValue("eventDate") // eventDate as string

	// Convert schoolID to integer
	schoolID, err := strconv.Atoi(schoolIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid school ID",
		})
	}

	// Handle file upload for heroImage if provided
	var heroImagePath string
	fileImage, err := c.FormFile("heroImage")
	if err == nil { // Process if a file is provided
		heroImagePath, err = utilities.UploadImage(fileImage, c, "announcement_image", "heroImage")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Create the request object for service layer
	createRequest := &request.AnnouncementCreateRequest{
		SchoolID:    schoolID,
		HeroImage:   heroImagePath,
		Title:       title,
		Description: description,
		Type:        announcementType,
		EventDate:   eventDateStr, // Set eventDate
	}

	// Call service to create announcement
	newAnnouncement, err := announcementController.announcementService.CreateAnnouncement(createRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Announcement created successfully.",
		"data":    newAnnouncement,
	})
}

// @Summary Get Announcement Types
// @Description Retrieve a list of available announcement types
// @Tags Announcements
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Success 200 {array} string "List of announcement types"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve announcement types"
// @Router /api/v1/announcement/type [get]
func (announcementController *AnnouncementController) GetAnnouncementTypes(c *fiber.Ctx) error {

	announcementTypes, err := announcementController.announcementService.GetAnnouncementTypes()
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve announcement type data",
		})
	}

	return c.JSON(fiber.Map{
		"data": announcementTypes,
	})
}

// @Summary Delete an Announcement
// @Description Delete an announcement by its ID
// @Tags Announcements
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "Announcement ID"
// @Success 200 {object} map[string]interface{} "Announcement deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Announcement not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete announcement"
// @Router /api/v1/announcement/delete/{id} [delete]
func (announcementController *AnnouncementController) DeleteAnnouncement(c *fiber.Ctx) error {
	err := utilities.CheckAccessAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	err = announcementController.announcementService.DeleteAnnouncement(uint(id), userID)
	if err != nil {
		if err.Error() == "announcement not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": constants.DataNotFoundMessage,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}

// @Summary Update an Announcement
// @Description Update an existing announcement by its ID
// @Tags Announcements
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "Announcement ID"
// @Param title formData string false "Title of the announcement"
// @Param description formData string false "Description of the announcement"
// @Param type formData string false "Type of the announcement (e.g., Event, General)"
// @Param eventDate formData string false "Event date in YYYY-MM-DD format"
// @Param heroImage formData file false "Hero image for the announcement"
// @Success 200 {object} map[string]interface{} "Announcement updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/announcement/update/{id} [put]
func (announcementController *AnnouncementController) UpdateAnnouncement(c *fiber.Ctx) error {
	// Parse announcement ID from URL
	id := c.Params("id")
	announcementID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid announcement ID",
		})
	}

	// Extract user ID from token claims
	var userID int
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	// Parse form fields
	title := c.FormValue("title")
	description := c.FormValue("description")
	announcementType := c.FormValue("type")
	eventDateStr := c.FormValue("eventDate")

	// Handle file upload for heroImage if provided
	var heroImage string
	fileHero, err := c.FormFile("heroImage")
	if err == nil { // Process if a file is provided
		heroImage, err = utilities.UploadImage(fileHero, c, "announcement_image", "heroImage")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Create the request object for service layer
	updateRequest := &request.AnnouncementUpdateRequest{
		Title:       title,
		Description: description,
		Type:        announcementType,
		EventDate:   eventDateStr,
		HeroImage:   heroImage,
	}

	// Call service to update announcement
	updatedAnnouncement, err := announcementController.announcementService.UpdateAnnouncement(announcementID, updateRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Announcement updated successfully.",
		"data":    updatedAnnouncement,
	})
}

// @Summary Get List of Announcements
// @Description Retrieve a paginated list of announcements
// @Tags Announcements
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc or desc, default: asc)"
// @Param type query string false "Announcement type filter"
// @Success 200 {object} map[string]interface{} "Paginated list of announcements"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve announcements"
// @Router /api/v1/announcement/getList [get]
func (ac *AnnouncementController) GetAnnouncementList(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "")
	sortOrder := c.Query("sortOrder", "asc")
	announcementType := c.Query("type", "")

	responseDTO, err := ac.announcementService.GetListAnnouncement(page, limit, search, sortBy, sortOrder, announcementType, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(responseDTO)
}

// @Summary Get Announcement Detail
// @Description Retrieve details of an announcement by its ID
// @Tags Announcements
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "Announcement ID"
// @Success 200 {object} map[string]interface{} "Announcement details"
// @Failure 400 {object} map[string]interface{} "Invalid announcement ID"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve announcement"
// @Router /api/v1/announcement/detail/{id} [get]
func (ac *AnnouncementController) GetAnnouncementDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	announcement, err := ac.announcementService.GetAnnouncementByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve announcement"})
	}
	return c.Status(200).JSON(announcement)
}
