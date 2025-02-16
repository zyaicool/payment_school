package controllers

import (
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BillingTypeController struct {
	billingTypeService services.BillingTypeService
}

func NewBillingTypeController(billingTypeService services.BillingTypeService) *BillingTypeController {
	return &BillingTypeController{billingTypeService: billingTypeService}
}

func (billingTypeController *BillingTypeController) GetAllBillingType(c *fiber.Ctx) error {
	// err := utilities.CheckAccessToBillingType(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")

	listBillingType, err := billingTypeController.billingTypeService.GetAllBillingType(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(listBillingType)
}

func (billingTypeController *BillingTypeController) GetDataBillingType(c *fiber.Ctx) error {
	err := utilities.CheckAccessToBillingType(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, _ := c.ParamsInt("id")
	billingType, err := billingTypeController.billingTypeService.GetBillingTypeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(billingType)
}

func (billingTypeController *BillingTypeController) CreateBillingType(c *fiber.Ctx) error {
	err := utilities.CheckAccessToBillingType(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var billingTypeRequest *request.BillingTypeCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&billingTypeRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createBillingType, err := billingTypeController.billingTypeService.CreateBillingType(billingTypeRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createBillingType,
	})
}

func (billingTypeController *BillingTypeController) UpdateBillingType(c *fiber.Ctx) error {
	err := utilities.CheckAccessToBillingType(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var billingTypeRequest *request.BillingTypeCreateUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&billingTypeRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedBillingType, err := billingTypeController.billingTypeService.UpdateBillingType(uint(id), billingTypeRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedBillingType,
	})
}

func (billingTypeController *BillingTypeController) DeleteBillingType(c *fiber.Ctx) error {
	err := utilities.CheckAccessToBillingType(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = billingTypeController.billingTypeService.DeleteBillingType(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}
