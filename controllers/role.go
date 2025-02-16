package controllers

import (
	"schoolPayment/constants"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type RoleController struct {
	roleService services.RoleService
}

func NewRoleController(roleService services.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

// @Summary Get All Roles
// @Description Retrieve all roles with pagination and search functionality
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit number of results" default(10)
// @Param search query string false "Search keyword"
// @Success 200 {object} []models.Role
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/role/getAllRole [get]
func (roleController *RoleController) GetAllRole(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")

	roles, err := roleController.roleService.GetAllRoles(page, limit, search, roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(roles)
}

// @Summary Get Data Role By Id
// @Description Get Data Role By Id
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/role/detail/{id} [get]
func GetDataRole(c *fiber.Ctx) error {
	// err := utilities.CheckAccessExceptUserOrtu(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	id, _ := c.ParamsInt("id")
	role, err := services.GetRoleByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(role)
}
