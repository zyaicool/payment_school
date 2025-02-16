package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupPaymentReportRoutes(api fiber.Router, paymentReportController *controllers.PaymentReportController) {
	apiPaymentReport := api.Group("/paymentReport")
	apiPaymentReport.Get("/getList", utilities.JWTProtected, paymentReportController.GetPaymentReport)
	apiPaymentReport.Get("/exportExcel", utilities.JWTProtected, paymentReportController.ExportPaymentReportToExcel)

}
