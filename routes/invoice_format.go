package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupInvoiceFormatRoutes(api fiber.Router, invoiceFormatController *controllers.InvoiceFormatController) {
	apiInvoiceFormat := api.Group("/invoiceFormat")
	apiInvoiceFormat.Post("/create", utilities.JWTProtected, invoiceFormatController.Create)
	apiInvoiceFormat.Get("/detail", utilities.JWTProtected, invoiceFormatController.GetBySchoolID)
}
