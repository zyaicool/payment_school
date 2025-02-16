package controllers

import (
	"schoolPayment/constants"
	"schoolPayment/models"
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type StudentParentController struct {
	studentParentService services.StudentParentService
}

func NewStudentParentController(studentParentService services.StudentParentService) *StudentParentController {
	return &StudentParentController{studentParentService: studentParentService}
}

func (studentParentController *StudentParentController) GetAllParent(c *fiber.Ctx) error {

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")

	parents, err := studentParentController.studentParentService.GetAllParent(page, limit, search, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(parents)
}

func (studentParentController *StudentParentController) GetDataParent(c *fiber.Ctx) error {

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	id, _ := c.ParamsInt("id")
	parent, err := studentParentController.studentParentService.GetParentByID(uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(parent)
}

func (studentParentController *StudentParentController) CreateDataParent(c *fiber.Ctx) error {

	var parent *models.StudentParent

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&parent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createParent, err := studentParentController.studentParentService.CreateParent(parent, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createParent,
	})
}

func (studentParentController *StudentParentController) UpdateParent(c *fiber.Ctx) error {

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	id, _ := c.ParamsInt("id")
	var parent *models.StudentParent

	if err := c.BodyParser(&parent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedParent, err := studentParentController.studentParentService.UpdateParent(uint(id), parent, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedParent,
	})
}

func (studentParentController *StudentParentController) GetDataParentByUserLogin(c *fiber.Ctx) error {

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	parent, err := studentParentController.studentParentService.GetParentByUserLogin(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(parent)
}

func (studentParentController *StudentParentController) CreateBatchDataParent(c *fiber.Ctx) error {

	var parents []*models.StudentParent

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&parents); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createParents, err := studentParentController.studentParentService.CreateBatchParent(parents, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createParents,
	})
}

func (studentParentController *StudentParentController) UpdateBatchParent(c *fiber.Ctx) error {

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	id, _ := c.ParamsInt("id")
	var parents []*models.StudentParent

	if err := c.BodyParser(&parents); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedParent, err := studentParentController.studentParentService.UpdateBatchParent(uint(id), parents, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedParent,
	})
}
