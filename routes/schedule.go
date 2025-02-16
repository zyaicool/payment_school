package routes

import (
	controllers "schoolPayment/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupScheduleRoutes(api fiber.Router, scheduleController *controllers.ScheduleController) {
	apiSchedule := api.Group("/schedule")
	apiSchedule.Get("/checkFailedTransaction", scheduleController.GetCheckPaymentFailedUsingSchedule)
	apiSchedule.Get("/reminderBilling", scheduleController.SendBillingReminder)
	apiSchedule.Get("/sendDummyNotif", scheduleController.SendDummyNotif)
}
