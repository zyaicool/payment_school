package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(api fiber.Router, dashController *controllers.DashboardController) {
	apiDashboard := api.Group("/dashboard")
	apiDashboard.Get("/parent", utilities.JWTProtected, dashController.GetDashboardFromParent)
	apiDashboard.Get("/admin", utilities.JWTProtected, dashController.GetDashboardFromAdmin)
}
