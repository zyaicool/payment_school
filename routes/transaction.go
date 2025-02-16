package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(api fiber.Router, transactionController *controllers.TransactionController) {
	apiTransaction := api.Group("/transaction")
	apiTransaction.Post("/create", utilities.JWTProtected, transactionController.CreateTransaction)
	apiTransaction.Post("/donation", utilities.JWTProtected, transactionController.PaymentDonation)

	apiMidtrans := apiTransaction.Group("/midtrans")
	apiMidtrans.Post("/payment", utilities.JWTProtected, transactionController.MidtransPayment)
	apiMidtrans.Get("/checkPayment", utilities.JWTProtected, controllers.MidtransCheckPayment)

	api.Post("/webhook", transactionController.HandleWebhook)
}
