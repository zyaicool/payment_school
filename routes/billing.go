package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupBillingRoutes(api fiber.Router, billingController *controllers.BillingController) {
	apiBilling := api.Group("/billing")
	apiBilling.Get("/getAllBilling", utilities.JWTProtected, billingController.GetAllBilling)
	apiBilling.Get("/detail/:id", utilities.JWTProtected, billingController.GetDataBilling)
	apiBilling.Post("/create", utilities.JWTProtected, billingController.CreateBilling)
	apiBilling.Put("/update/:id", utilities.JWTProtected, controllers.UpdateBilling)
	apiBilling.Delete("/delete/:id", utilities.JWTProtected, controllers.DeleteBilling)
	apiBilling.Get("/generateInstallment/:id", utilities.JWTProtected, controllers.GenerateInstallment)
	apiBilling.Get("/billingStatus", billingController.GetBillingStatuses)
	apiBilling.Post("/createDonation", utilities.JWTProtected, billingController.CreateDonation)
	apiBilling.Get("/getBillingByStudentID", utilities.JWTProtected, billingController.GetBillingByStudentID)
}
