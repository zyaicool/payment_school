package routes

import (
	// "os"
	"net/http"

	controllers "schoolPayment/controllers"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app fiber.Router) {
	provinceService := services.NewProvinceService()
	cityService := services.NewCityService()
	paymentTypeService := services.NewPaymentTypeService()
	paymentStatusService := services.NewPaymentStatusService()
	paymentStatusController := controllers.NewPaymentController(paymentStatusService)

	provinceController := controllers.NewProvinceController(provinceService)
	cityController := controllers.NewCityController(cityService)
	paymentTypeController := controllers.NewPaymentTypeController(paymentTypeService)

	app.Static("/upload", "./upload")
	app.Get("/province", provinceController.GetDataProvinces)
	app.Get("/city", cityController.GetDataCities)
	app.Get("/paymentType", paymentTypeController.ListPaymentTypes)
	app.Get("/paymentStatus", paymentStatusController.GetListPaymentStatus)
	app.Post("/generateMasterData", controllers.GenerateMasterData)
	app.Get("/health", func(c *fiber.Ctx) error {
		// Simulating a health check
		isHealthy := true // Change this to false to simulate an error

		if isHealthy {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"status": "healthy",
			})
		} else {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  "Service is not running properly",
			})
		}
	})

}
