package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router, roleController *controllers.RoleController) {
	apiRole := api.Group("/role")
	apiRole.Get("/getAllRole", utilities.JWTProtected, roleController.GetAllRole)
	apiRole.Get("/detail/:id", utilities.JWTProtected, controllers.GetDataRole)
}
