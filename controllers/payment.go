package controllers

import (
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	paymentService services.PaymentStatusService
}

func NewPaymentController(paymentService services.PaymentStatusService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

// @Summary Get List of Payment Statuses
// @Description Retrieve a list of available payment statuses from a JSON data source
// @Tags Payment
// @Accept json
// @Produce json
// @Success 200 {object} []response.PaymentType "List of payment statuses"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve payment statuses"
// @Router /api/v1/paymentStatus [get]
func (controller *PaymentController) GetListPaymentStatus(c *fiber.Ctx) error {
	paymentStatuses, err := controller.paymentService.GetPaymentStatus("data/payment_status.json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": paymentStatuses,
	})
}
