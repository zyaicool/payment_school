package routes

import (
	controllers "schoolPayment/controllers"
	utilities "schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
)

func SetupBankAccountRoutes(api fiber.Router, bankAccountController *controllers.BankAccountController) {
	apiBankAccount := api.Group("/bankAccount") 
	apiBankAccount.Get("/getListBankAccount", utilities.JWTProtected, bankAccountController.GetAllBankAccounts)
	apiBankAccount.Get("/detail/:id", utilities.JWTProtected, bankAccountController.GetBankAccountByID)
	apiBankAccount.Post("/create", utilities.JWTProtected, bankAccountController.CreateBankAccount)
	apiBankAccount.Put("/update/:id", utilities.JWTProtected, bankAccountController.UpdateBankAccount)
	apiBankAccount.Delete("delete/:id", utilities.JWTProtected, bankAccountController.DeleteBankAccount)
	apiBankAccount.Get("/listBankName", utilities.JWTProtected, bankAccountController.GetBankName)
}
