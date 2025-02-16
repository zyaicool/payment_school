package controllers

import (
	"schoolPayment/constants"
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PrefixClassServiceInterface interface {
	GetAllPrefixClassService(search string, userID int) (response.PrefixCLassResponse, error)
	CreatePrefixClassService(prefixClass *request.PrefixClassCreate, userID int) (*models.PrefixClass, error)
}

// Update the controller struct
type PrefixClassController struct {
	prefixClassService PrefixClassServiceInterface
}

func NewPrefixClassController() *PrefixClassController {
	prefixClassService := services.NewPrefixClassService()
	return &PrefixClassController{
		prefixClassService: prefixClassService,
	}
}

// @Summary Get Prefix Class
// @Description Retrieve a list of prefix classes filtered by search criteria
// @Tags PrefixClass
// @Accept json
// @Produce json
// @Param search query string false "Search filter for prefix classes"
// @Success 200 {array} response.PrefixCLassResponse "List of prefix classes"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve prefix classes"
// @Router /api/v1/prefixClass/getListPrefixClass [get]
func (prefixClassController *PrefixClassController) GetPrefixClass(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["user_id"].(float64))

	search := c.Query("search")

	prefix, err := prefixClassController.prefixClassService.GetAllPrefixClassService(search, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(prefix)
}

// @Summary Create a Prefix Class
// @Description Create a new prefix class with the provided details
// @Tags PrefixClass
// @Accept json
// @Produce json
// @Param prefixClass body request.PrefixClassCreate true "Prefix Class Data"
// @Success 200 {object} models.PrefixClass "Prefix class created successfully"
// @Failure 400 {object} map[string]interface{} "Validation error or missing fields"
// @Failure 500 {object} map[string]interface{} "Failed to create prefix class"
// @Router /api/v1/prefixClass/create [post]
func (prefixClassController *PrefixClassController) CreatePrefixClass(c *fiber.Ctx) error {
	var prefixClass *request.PrefixClassCreate
	var userID int = 0

	// Get user information from the token claims stored in context
	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		// If the user claims exist and are of type jwt.MapClaims, extract the user ID.
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	if err := c.BodyParser(&prefixClass); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	if prefixClass.PrefixName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "PrefixName is required",
		})
	}
	if prefixClass.SchoolID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SchoolID is required",
		})
	}

	createPrefix, err := prefixClassController.prefixClassService.CreatePrefixClassService(prefixClass, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Data berhasil disimpan.",
		"data":    createPrefix,
	})
}
