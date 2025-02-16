package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupSchoolMajorRoutes(api fiber.Router, schoolMajorController *controllers.SchoolMajorController) {
	apiSchoolMajor := api.Group("/schoolMajor") 
	apiSchoolMajor.Get("/getListSchoolMajor", utilities.JWTProtected, schoolMajorController.GetAllSchoolMajor)
	apiSchoolMajor.Post("/create", utilities.JWTProtected, schoolMajorController.CreateSchoolMajor)
}
