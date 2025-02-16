package controllers

import (
	"schoolPayment/constants"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Update the controller to use the interface instead of the concrete type
type BillingHistoryController struct {
	billingHistoryService services.BillingHistoryServiceInterface
}

func NewBillingHistoryController(billingHistoryService services.BillingHistoryServiceInterface) *BillingHistoryController {
	return &BillingHistoryController{billingHistoryService: billingHistoryService}
}

// @Summary Get All Billing History
// @Description Retrieve all billing history based on various filters and pagination
// @Tags Billing History
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Param search query string false "Search keyword"
// @Param studentId query int false "Student ID filter"
// @Param paymentStatusCode query string false "Payment status code filter"
// @Param schoolYearId query int false "School year ID filter"
// @Param paymentTypeId query int false "Payment type ID filter"
// @Param userId query int false "User ID filter"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc or desc)"
// @Success 200 {array} response.ListBillingHistory "List of billing history records"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve billing history"
// @Router /api/v1/billingHistory/getAllBillingHistory [get]
func (billingHistoryController BillingHistoryController) GetAllBillingHistory(c *fiber.Ctx) error {
	// err := utilities.CheckAccessExceptUserOrtu(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	userLoginID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")
	studentID := c.QueryInt("studentId", 0)
	paymentStatusCode := c.Query("paymentStatusCode")
	schoolYearID := c.QueryInt("schoolYearId", 0)
	paymentTypeID := c.QueryInt("paymentTypeId")
	userID := c.QueryInt("userId", 0)
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")

	listBilling, err := billingHistoryController.billingHistoryService.GetAllBillingHistory(page, limit, search, studentID, roleID, schoolYearID, paymentTypeID, paymentStatusCode, userID, userLoginID, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(listBilling)
}

// @Summary Get Billing History Details
// @Description Retrieve detailed billing history for a specific transaction ID
// @Tags Billing History
// @Accept json
// @Produce json
// @Param transactionId query int true "Transaction ID"
// @Success 200 {object} response.DetailBillingHistory "Billing history details"
// @Failure 404 {object} map[string]interface{} constants.DataNotFoundMessage
// @Router /api/v1/billingHistory/detailBillingHistory [get]
func (billingHistoryController BillingHistoryController) GetDetailBillingHistoryID(c *fiber.Ctx) error {
	// err := utilities.CheckAccessExceptUserOrtu(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	transactionId := c.QueryInt("transactionId")

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	billing, err := billingHistoryController.billingHistoryService.GetDetailBillingHistoryIDService(transactionId, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(billing)
}

// @Summary Generate Invoice PDF
// @Description Generate a PDF invoice for a given invoice number
// @Tags Billing History
// @Accept json
// @Produce application/pdf
// @Param invoiceNumber query string true "Invoice number"
// @Success 200 {file} file "PDF file generated successfully"
// @Failure 500 {object} map[string]interface{} "Failed to generate invoice"
// @Router /api/v1/billingHistory/printInvoice [get]
func (billingHistoryController BillingHistoryController) GeneratePDF(c *fiber.Ctx) error {
	invoiceNumber := c.Query("invoiceNumber")

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	_, err := billingHistoryController.billingHistoryService.GenerateInvoice(c, invoiceNumber, true, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}
