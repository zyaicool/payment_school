package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupBillingStudentRoutes(api fiber.Router, billingStudentController *controllers.BillingStudentController) {
	apiBillingStudent := api.Group("/billingStudent")
	apiBillingStudent.Get("/getAllBillingStudent", utilities.JWTProtected, billingStudentController.GetAllBillingStudent)
	apiBillingStudent.Get("/detailBillingStudent", utilities.JWTProtected, billingStudentController.GetDetailBillingID)
	apiBillingStudent.Get("/detailBillingByBillingStudentId/:id", utilities.JWTProtected, billingStudentController.GetDetailBillingStudent)
	apiBillingStudent.Put("/update/:id", utilities.JWTProtected, billingStudentController.UpdateBillingStudent)
	apiBillingStudent.Delete("/delete/:id", utilities.JWTProtected, billingStudentController.DeleteBillingStudent)
	apiBillingStudent.Post("/create", utilities.JWTProtected, billingStudentController.CreateBillingStudent)

}
