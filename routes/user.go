package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, userController *controllers.UserController) {
	apiUser := api.Group("/user") // Group routes under "/v1/user"
	apiUser.Get("/getAllUser", utilities.JWTProtected, userController.GetAllUser)
	apiUser.Get("/detail/:id", utilities.JWTProtected, userController.GetDataUser)
	apiUser.Get("/getByEmail", controllers.GetDataByEmail)
	apiUser.Get("/generateTokenChangePassword", userController.GenerateTokenChangePassword)
	apiUser.Get("/verifyTokenChangePassword", userController.VerifyTokenChangePassword)
	apiUser.Post("/create", utilities.JWTProtected, userController.CreateUser)
	apiUser.Post("/selfRegistration", userController.CreateUser)
	apiUser.Put("/edit/:id", utilities.JWTProtected, userController.UpdateUser)
	apiUser.Put("/verifyEmail", userController.EmailVerification)
	apiUser.Put("/changePassword", userController.ChangePassword)
	apiUser.Delete("/delete/:id", utilities.JWTProtected, userController.DeleteUser)
	apiUser.Get("/getFileExcel", utilities.JWTProtected, userController.DownloadFileExcelFormatForUser)
	apiUser.Post("/upload", utilities.JWTProtected, userController.UploadUsers)
	apiUser.Get("/checkEmail", userController.CheckExistingEmail)
	apiUser.Get("/checkUsername", controllers.CheckExistingUsername)
	apiUser.Get("/resendVerificationEmail/:id", userController.ResendEmailVerification)
	apiUser.Put("/uploadFoto", utilities.JWTProtected, userController.UploadUserPhoto)
	apiUser.Put("/updatePassword/:id", utilities.JWTProtected, userController.ChangePasswordWithoutToken)	
}
