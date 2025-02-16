package controllers

import (
	"schoolPayment/dtos/request"
	"schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

type ScheduleController struct {
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) *ScheduleController {
	return &ScheduleController{scheduleService: scheduleService}
}

// @Summary Check for Payment Failures Using Schedule
// @Description Retrieves billing transactions that have failed payments based on the schedule
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of billing transactions with payment failure"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve billing transactions"
// @Router /api/v1/schedule/checkFailedTransaction [get]
func (scheduleController ScheduleController) GetCheckPaymentFailedUsingSchedule(c *fiber.Ctx) error {
	// userClaims := c.Locals("user").(jwt.MapClaims)
	// userID := int(userClaims["user_id"].(float64))

	billingTransaction, err := scheduleController.scheduleService.GetCheckPaymentFailedUsingScheduleService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(billingTransaction)
}

// @Summary Send a Dummy Notification
// @Description Sends a dummy notification for testing purposes
// @Tags Schedule
// @Accept json
// @Produce json
// @Param body body request.DummyNotifRequest true "Notification details"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 500 {object} map[string]interface{} "Failed to send notification"
// @Router /api/v1/schedule/sendDummyNotif [post]
func (scheduleController ScheduleController) SendDummyNotif(c *fiber.Ctx) error {
	var schedule *request.DummyNotifRequest
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	err :=  scheduleController.scheduleService.DummySendNotif(5, schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Send Dummy Notif",
	})
}


func  (scheduleController ScheduleController) SendBillingReminder(c *fiber.Ctx) error {
	// userClaims := c.Locals("user").(jwt.MapClaims)
	// userID := int(userClaims["user_id"].(float64))

	succses, err := scheduleController.scheduleService.SendReminderDueDate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(succses)
}
