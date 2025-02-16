package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupLoginRoutes(api fiber.Router, loginController *controllers.LoginController) {
	api.Post("/login", loginController.Login)
	api.Get("/auth", utilities.JWTProtected, loginController.AuthUser)
	api.Post("/invalidFirebaseToken",utilities.JWTProtected, loginController.InvalidateFirebaseToken)
}
