package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupBillingHistoryRoutes(api fiber.Router, bilingHistoryController *controllers.BillingHistoryController) {
	apiBillingHistory := api.Group("/billingHistory") 
	apiBillingHistory.Get("/getAllBillingHistory", utilities.JWTProtected, bilingHistoryController.GetAllBillingHistory)
	apiBillingHistory.Get("/detailBillingHistory", utilities.JWTProtected, bilingHistoryController.GetDetailBillingHistoryID)
	apiBillingHistory.Get("/printInvoice", utilities.JWTProtected, bilingHistoryController.GeneratePDF)
}
