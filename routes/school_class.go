package routes

import (
	controllers "schoolPayment/controllers"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupSchoolClassRoutes(api fiber.Router) {
	// Initialize controller directly
	schoolClassController := controllers.NewSchoolClassController()

	// Setup routes
	apiSchoolClass := api.Group("/schoolClass")
	apiSchoolClass.Get("/getAllSchoolClass", utilities.JWTProtected, schoolClassController.GetAllSchoolClass)
	apiSchoolClass.Get("/detail/:id", utilities.JWTProtected, schoolClassController.GetDataSchoolClass)
	apiSchoolClass.Post("/create", utilities.JWTProtected, schoolClassController.CreateSchoolClass)
	apiSchoolClass.Put("/update/:id", utilities.JWTProtected, schoolClassController.UpdateSchoolClass)
	apiSchoolClass.Delete("/delete/:id", utilities.JWTProtected, schoolClassController.DeleteSchoolClass)
}
