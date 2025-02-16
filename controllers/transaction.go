package controllers

import (
	"fmt"

	"schoolPayment/constants"
	"schoolPayment/dtos/request"
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

func (transactionController *TransactionController) GetAllTransaction(c *fiber.Ctx) error {

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// get query params pagination and email
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")

	transactions, err := transactionController.transactionService.GetAllTransaction(page, limit, search, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(transactions)
}

// MidtransPayment processes the Midtrans payment for a transaction.
// @Summary Make Midtrans Payment
// @Description Process a payment using Midtrans payment gateway.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param request body request.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/transaction/midtrans/payment [post]
func (transactionController *TransactionController) MidtransPayment(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	var request request.CreateTransactionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.FailedToParseRequestBodyMessage,
		})
	}

	rsp, err := transactionController.transactionService.MidtransPayment(
		&services.BillingService{},
		request.StudentId,
		request.PaymentMethodId,
		userID,
		request.BillingStudentIds,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rsp)
}

// PaymentDonation processes a donation payment.
// @Summary Make Payment Donation
// @Description Process a payment donation for a student.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param request body request.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/transaction/donation [post]
func (transactionController *TransactionController) PaymentDonation(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	var request request.CreateTransactionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.FailedToParseRequestBodyMessage,
		})
	}

	rsp, err := transactionController.transactionService.PaymentDonation(
		request.StudentId,
		request.BillingId,
		request.Amount,
		request.PaymentMethodId,
		userID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rsp)
}

// MidtransCheckPayment checks the status of a Midtrans payment.
// @Summary Check Midtrans Payment Status
// @Description Check the payment status of a Midtrans transaction by order ID.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param orderId query string true "Midtrans Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/transaction/midtrans/checkPayment [get]
func MidtransCheckPayment(c *fiber.Ctx) error {
	// userClaims := c.Locals("user").(jwt.MapClaims)
	// userID := int(userClaims["user_id"].(float64))

	orderID := c.Query("orderId")

	resp, err := services.MidtransCheckTransaction(orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(resp)
}

// CreateTransaction creates a new transaction.
// @Summary Create New Transaction
// @Description Create a new transaction for a student with the provided payment method.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param request body request.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/transaction/create [post]
func (transactionController *TransactionController) CreateTransaction(c *fiber.Ctx) error {
	var request request.CreateTransactionRequest
	var userId int = 0

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.FailedToParseRequestBodyMessage,
		})
	}

	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userId = int(userClaimID)
		}
	}

	CreateTransaction, err := transactionController.transactionService.CreateTransactionService(request, userId, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    CreateTransaction,
	})
}

// HandleWebhook processes a webhook from payment gateway.
// @Summary Handle Payment Webhook
// @Description Handle incoming webhook from the payment gateway for transaction status updates.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param payload body request.WebhookPayload true "Webhook Payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/webhook [post]
func (transactionController *TransactionController) HandleWebhook(c *fiber.Ctx) error {
	// Parse incoming JSON payload
	var payload request.WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse JSON",
		})
	}

	// Check if SettlementTime is available
	settlementTime := "N/A"
	if payload.SettlementTime != nil {
		settlementTime = *payload.SettlementTime
	}

	// Log or process the webhook event
	fmt.Printf("Received webhook: Transaction ID: %s, Status: %s, Settlement Time: %s, Gross Amount: %s\n",
		payload.TransactionID, payload.TransactionStatus, settlementTime, payload.GrossAmount)

	err := transactionController.transactionService.UpdateFromWebHook(payload, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Respond with a success message
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Webhook received",
	})
}
