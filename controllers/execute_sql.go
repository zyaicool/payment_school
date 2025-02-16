package controllers

import (
	config "schoolPayment/configs"
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
)

func GenerateMasterData(c *fiber.Ctx) error {
	// execute sql master data
	sqlFilePath := "./Script-Insert-Master-Data.sql"
	err := config.ExecuteSQLScriptUsingGORM(config.DB, sqlFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = services.CreateDataHistoryMasterData(sqlFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
