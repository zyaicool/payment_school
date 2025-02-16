package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupSchoolGradeRoutes(api fiber.Router, schoolGradeController *controllers.SchoolGradeController) {
	apiSchoolGrade := api.Group("/schoolGrade") 
	apiSchoolGrade.Get("/getAllSchoolGrade", utilities.JWTProtected, controllers.GetAllSchoolGrade)
	apiSchoolGrade.Get("/detail/:id", utilities.JWTProtected, schoolGradeController.GetDataSchoolGrade)
	apiSchoolGrade.Post("/create", utilities.JWTProtected, schoolGradeController.CreateSchoolGrade)
	apiSchoolGrade.Put("/update/:id", utilities.JWTProtected, schoolGradeController.UpdateSchoolGrade)
	apiSchoolGrade.Delete("/delete/:id", utilities.JWTProtected, schoolGradeController.DeleteSchoolGrade)
}
