package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupBillingReportRoutes(api fiber.Router, billingReportController *controllers.BillingReportController) {
	apiBillingReport := api.Group("/billingReport")
	apiBillingReport.Get("/getList", utilities.JWTProtected, billingReportController.GetBillingReports)
	apiBillingReport.Get("/exportExcel", utilities.JWTProtected, billingReportController.ExportBillingReportToExcel)
}
