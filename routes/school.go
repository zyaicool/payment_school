package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupSchoolRoutes(api fiber.Router, schoolController *controllers.SchoolController) {
	apiSchool := api.Group("/school")
	apiSchool.Get("/getAllSchool", utilities.JWTProtected, schoolController.GetAllSchool)
	apiSchool.Get("/detail/:id", utilities.JWTProtected, schoolController.GetDataSchool)
	apiSchool.Post("/create", utilities.JWTProtected, schoolController.CreateSchool)
	apiSchool.Put("/update/:id", utilities.JWTProtected, schoolController.UpdateSchool)
	apiSchool.Delete("/delete/:id", utilities.JWTProtected, schoolController.DeleteSchool)
	apiSchool.Get("/getAllOnboarding", schoolController.GetAllOnboardingSchools)

}
