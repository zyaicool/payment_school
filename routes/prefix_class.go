package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupPrefixClassRoutes(api fiber.Router) {
	prefixClassController := controllers.NewPrefixClassController()

	apiPrefixClass := api.Group("/prefixClass")
	apiPrefixClass.Get("/getListPrefixClass", utilities.JWTProtected, prefixClassController.GetPrefixClass)
	apiPrefixClass.Post("/create", utilities.JWTProtected, prefixClassController.CreatePrefixClass)
}
