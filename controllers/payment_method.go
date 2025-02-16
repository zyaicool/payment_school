package controllers

import (
	"strconv"

	request "schoolPayment/dtos/request"
	"schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PaymentMethodController struct {
	paymentMethodService services.PaymentMethodService
}

func NewPaymentMethodController(paymentMethodService services.PaymentMethodService) PaymentMethodController {
	return PaymentMethodController{paymentMethodService: paymentMethodService}
}

// @Summary Create Payment Method
// @Description Add a new payment method with optional logo upload
// @Tags Payment Method
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param paymentMethod formData string true "Payment Method Name"
// @Param bankCode formData string true "Bank Code"
// @Param bankName formData string true "Bank Name"
// @Param adminFee formData int true "Administrative Fee"
// @Param isPercentage formData string false "Is Percentage (true/false)"
// @Param adminFeePercentage formData string false "Administrative Fee Percentage"
// @Param methodLogo formData file false "Logo for the Payment Method"
// @Success 200 {object} models.PaymentMethod "Data saved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/masterPaymentGateway/create [post]
func (paymentMethodController *PaymentMethodController) CreatePaymentMethod(c *fiber.Ctx) error {
	var userID int = 0

	// Extract user ID from token claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	// Parse form fields
	paymentMethod := c.FormValue("paymentMethod")
	bankCode := c.FormValue("bankCode")
	bankName := c.FormValue("bankName")
	adminFeeStr := c.FormValue("adminFee")
	isPercentage := c.FormValue("isPercentage", "false")
	adminFeePercentage := c.FormValue("adminFeePercentage")

	// Convert adminFee to integer
	adminFee, err := strconv.Atoi(adminFeeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid admin fee value",
		})
	}

	// Handle file upload for methodLogo if provided
	var methodLogo string
	fileLogo, err := c.FormFile("methodLogo")
	if err == nil { // Process if a file is provided
		methodLogo, err = utilities.UploadImage(fileLogo, c, "payment_method_logo", "methodLogo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Create the request object for service layer
	createRequest := &request.PaymentMethodCreateRequest{
		PaymentMethod:      paymentMethod,
		BankCode:           bankCode,
		BankName:           bankName,
		AdminFee:           adminFee,
		MethodLogo:         methodLogo,
		IsPercentage:       isPercentage,
		AdminFeePercentage: adminFeePercentage,
	}

	// Call service to create payment method
	newPaymentMethod, err := paymentMethodController.paymentMethodService.CreatePaymentMethod(createRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data saved successfully.",
		"data":    newPaymentMethod,
	})
}

// @Summary Update Payment Method
// @Description Update an existing payment method with optional logo upload
// @Tags Payment Method
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "Payment Method ID"
// @Param paymentMethod formData string true "Payment Method Name"
// @Param bankCode formData string true "Bank Code"
// @Param bankName formData string true "Bank Name"
// @Param adminFee formData int true "Administrative Fee"
// @Param isPercentage formData string false "Is Percentage (true/false)"
// @Param adminFeePercentage formData string false "Administrative Fee Percentage"
// @Param methodLogo formData file false "Logo for the Payment Method"
// @Success 200 {object} models.PaymentMethod "Payment method updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input or ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/masterPaymentGateway/update/{id} [put]
func (pmc *PaymentMethodController) UpdatePaymentMethod(c *fiber.Ctx) error {
	// Parse payment method ID from URL
	id := c.Params("id")
	paymentMethodID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment method ID",
		})
	}

	// Parse form fields
	paymentMethod := c.FormValue("paymentMethod")
	bankCode := c.FormValue("bankCode")
	bankName := c.FormValue("bankName")
	adminFeeStr := c.FormValue("adminFee")
	isPercentage := c.FormValue("isPercentage")
	adminFeePercentage := c.FormValue("adminFeePercentage")

	// Convert adminFee to integer
	adminFee, err := strconv.Atoi(adminFeeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid admin fee value",
		})
	}

	// Handle file upload for methodLogo if provided
	var methodLogo string
	fileLogo, err := c.FormFile("methodLogo")
	if err == nil { // Process if a file is provided
		methodLogo, err = utilities.UploadImage(fileLogo, c, "payment_method_logo", "methodLogo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Create the request object for service layer
	updateRequest := &request.PaymentMethodCreateRequest{
		PaymentMethod:      paymentMethod,
		BankCode:           bankCode,
		BankName:           bankName,
		AdminFee:           adminFee,
		MethodLogo:         methodLogo,
		IsPercentage:       isPercentage,
		AdminFeePercentage: adminFeePercentage,
	}

	// Call service to update payment method
	updatedPaymentMethod, err := pmc.paymentMethodService.UpdatePaymentMethod(paymentMethodID, updateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment method updated successfully.",
		"data":    updatedPaymentMethod,
	})
}

// @Summary Get All Payment Methods
// @Description Retrieve a list of all payment methods with optional search functionality
// @Tags Payment Method
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param search query string false "Search term for payment method"
// @Success 200 {object} []response.PaymentMethodResponse "List of payment methods"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/masterPaymentGateway/getAllConfig [get]
func (pmc *PaymentMethodController) GetAllPaymentMethod(c *fiber.Ctx) error {
	// Get search query parameter
	search := c.Query("search", "") // Default is an empty string if no search is provided

	// Call service to get the payment gateway config data
	data, err := pmc.paymentMethodService.GetAllPaymentMethod(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// @Summary Get Payment Method Detail
// @Description Retrieve detailed information of a payment method by its ID
// @Tags Payment Method
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "Payment Method ID"
// @Success 200 {object} response.PaymentMethodResponse "Payment method details"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/masterPaymentGateway/detail/{id} [get]
func (controller *PaymentMethodController) GetPaymentMethodDetail(c *fiber.Ctx) error {
	// Get ID from URL parameters
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}

	// Fetch the payment method detail
	response, err := controller.paymentMethodService.GetPaymentMethodDetail(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}
