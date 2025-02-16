package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupParentRoutes(api fiber.Router, studentParentController *controllers.StudentParentController) {
	apiParent := api.Group("/parent") // Group parent under "/v1/parent"
	apiParent.Get("/getAllParent", utilities.JWTProtected, studentParentController.GetAllParent)
	apiParent.Get("/detail/:id", utilities.JWTProtected, studentParentController.GetDataParent)
	apiParent.Post("/create", utilities.JWTProtected, studentParentController.CreateDataParent)
	apiParent.Post("/createBatch", utilities.JWTProtected, studentParentController.CreateBatchDataParent)
	apiParent.Put("/edit/:id", utilities.JWTProtected, studentParentController.UpdateParent)
	apiParent.Put("/editBatch", utilities.JWTProtected, studentParentController.UpdateParent)
	apiParent.Get("/getByUser", utilities.JWTProtected, studentParentController.GetDataParentByUserLogin)
}
