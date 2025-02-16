package controllers

import (
	"net/http"

	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

type PaymentTypeController struct {
	paymentTypeService services.PaymentTypeService
}

func NewPaymentTypeController(paymentTypeService services.PaymentTypeService) *PaymentTypeController {
	return &PaymentTypeController{paymentTypeService: paymentTypeService}
}

// @Summary List Payment Types
// @Description Retrieve a list of available payment types from a JSON data source
// @Tags Payment Type
// @Accept json
// @Produce json
// @Success 200 {object} []response.PaymentType "List of payment types"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve payment types"
// @Router /api/v1/paymentType [get]
func (PaymentTypeController *PaymentTypeController) ListPaymentTypes(ctx *fiber.Ctx) error {
	paymentTypes, err := PaymentTypeController.paymentTypeService.GetDataPaymentTypes("data/payment_type.json")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve payment types",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": paymentTypes,
	})
}
