package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupSchoolYearRoutes(api fiber.Router, schoolYearController *controllers.SchoolYearController) {
	apiSchoolYear := api.Group("/schoolYear")
	apiSchoolYear.Get("/getAllSchoolYear", utilities.JWTProtected, schoolYearController.GetAllSchoolYear)
	apiSchoolYear.Get("/detail/:id", utilities.JWTProtected, schoolYearController.GetDataSchoolYear)
	apiSchoolYear.Post("/create", utilities.JWTProtected, schoolYearController.CreateSchoolYear)
	apiSchoolYear.Put("/update/:id", utilities.JWTProtected, schoolYearController.UpdateSchoolYear)
	apiSchoolYear.Delete("/delete/:id", utilities.JWTProtected, schoolYearController.DeleteSchoolYear)
}
