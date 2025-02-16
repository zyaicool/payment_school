package routes

import (
	"schoolPayment/controllers"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router) {
	studentController := controllers.GetStudentController()
	apiStudent := api.Group("/student")

	apiStudent.Get("/getAllStudent", utilities.JWTProtected, studentController.GetAllStudent)
	apiStudent.Get("/detail/:id", utilities.JWTProtected, studentController.GetStudentByID)
	apiStudent.Post("/create", utilities.JWTProtected, studentController.CreateStudent)
	apiStudent.Put("/update/:id", utilities.JWTProtected, studentController.UpdateStudent)
	apiStudent.Delete("/delete/:id", utilities.JWTProtected, studentController.DeleteStudent)
	apiStudent.Post("/linkToUser", utilities.JWTProtected, studentController.CreateLinkToUser)
	apiStudent.Get("/user/:user_id", utilities.JWTProtected, studentController.GetStudentByUserId)
	apiStudent.Get("/getFileExcel", utilities.JWTProtected, studentController.DownloadFileExcelFormatForStudent)
	apiStudent.Post("/upload", utilities.JWTProtected, studentController.UploadStudents)
	apiStudent.Get("/image/:id", utilities.JWTProtected, studentController.GetStudentImage)
	apiStudent.Post("/image/:id", utilities.JWTProtected, studentController.UploadImageStudent)
	apiStudent.Get("/exportExcel", utilities.JWTProtected, studentController.ExportStudentToExcel)
}
