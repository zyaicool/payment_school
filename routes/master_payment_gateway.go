package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupPaymentMethodRoutes(api fiber.Router, paymentMethodController *controllers.PaymentMethodController) {
	apiPaymentMethod := api.Group("/masterPaymentGateway")
	apiPaymentMethod.Get("/getAllConfig", utilities.JWTProtected, paymentMethodController.GetAllPaymentMethod)
	apiPaymentMethod.Get("/detail/:id", utilities.JWTProtected, paymentMethodController.GetPaymentMethodDetail)
	apiPaymentMethod.Post("/create", utilities.JWTProtected, paymentMethodController.CreatePaymentMethod)
	apiPaymentMethod.Put("/update/:id", utilities.JWTProtected, paymentMethodController.UpdatePaymentMethod)
}
