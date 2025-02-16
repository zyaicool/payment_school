package controllers

import (
	"schoolPayment/dtos/request"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


type InvoiceFormatServiceInterface interface {
	Create(request *request.CreateInvoiceFormatRequest, userID int) (*models.InvoiceFormat, error)
	GetBySchoolID(schoolID uint) (*models.InvoiceFormat, error)
}

type InvoiceFormatController struct {
	service InvoiceFormatServiceInterface
}

func NewInvoiceFormatController(service InvoiceFormatServiceInterface) *InvoiceFormatController {
	return &InvoiceFormatController{
		service: service,
	}
}

// @Summary Create Invoice Format
// @Description Create a new invoice format for a school
// @Tags Invoice Format
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param body body request.CreateInvoiceFormatRequest true "Create Invoice Format Request"
// @Success 200 {object} models.InvoiceFormat "Invoice format created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/invoiceFormat/create [post]
func (c *InvoiceFormatController) Create(ctx *fiber.Ctx) error {
	var request request.CreateInvoiceFormatRequest

	// Get user information from the token claims stored in context
	var userID int = 0

	// Get user information from the token claims stored in context
	userClaims, ok := ctx.Locals("user").(jwt.MapClaims)

	if !ok || userClaims == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized or invalid token",
		})
	}
	
	if userIDValue, ok := userClaims["user_id"].(float64); ok {
		userID = int(userIDValue)
	} else {
		userID = 0
	}
	// Parse the request body
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the user has access to create the invoice format (admin school check)
	if err := utilities.CheckAccessAdminSekolah(ctx); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Call the service to create or update the invoice format
	invoiceFormat, err := c.service.Create(&request, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    invoiceFormat,
		"message": "Invoice format created successfully",
	})
}

// @Summary Get Invoice Format by School ID
// @Description Retrieve an invoice format by its associated school ID
// @Tags Invoice Format
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param schoolId query int true "School ID"
// @Success 200 {object} response.InvoiceFormatResponse "Invoice format data"
// @Failure 400 {object} map[string]interface{} "Invalid school ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/invoiceFormat/detail [get]
func (c *InvoiceFormatController) GetBySchoolID(ctx *fiber.Ctx) error {
	// Get the school ID from query params
	schoolID := ctx.QueryInt("schoolId")
	if schoolID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid school ID",
		})
	}

	// Call the service to get the invoice format by school ID
	invoiceFormat, err := c.service.GetBySchoolID(uint(schoolID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the invoice format in the response
	return ctx.JSON(response.InvoiceFormatResponse{
		ID:                     invoiceFormat.ID,
		SchoolID:               invoiceFormat.SchoolID,
		Prefix:                 invoiceFormat.Prefix,
		Format:                 invoiceFormat.Format,
		GeneratedInvoiceFormat: invoiceFormat.GeneratedInvoiceFormat,
		CreatedBy:              invoiceFormat.CreatedBy,
		UpdatedBy:              invoiceFormat.UpdatedBy,
		CreatedAt:              invoiceFormat.CreatedAt,
		UpdatedAt:              invoiceFormat.UpdatedAt,
	})
}
