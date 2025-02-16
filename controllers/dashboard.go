package controllers

import (
	"schoolPayment/models"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type DashboardController struct {
	dashboardService services.DashboardService
}

func NewDashboardController(dashboardService services.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

// @Summary Get Dashboard Data for Parent
// @Description Retrieve dashboard data for a parent user, optionally filtered by student ID
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param studentId query int false "Student ID"
// @Success 200 {object} response.DashboardResponse "Dashboard data for parent"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/dashboard/parent [get]
func (dashboardController *DashboardController) GetDashboardFromParent(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	roleID := int(userClaims["role_id"].(float64))
	schoolID := int(userClaims["school_id"].(float64))

	var user = models.User{
		Master: models.Master{
			ID: uint(userID),
		},
		RoleID: uint(roleID),
		UserSchool: &models.UserSchool{
			SchoolID: uint(schoolID),
			School: &models.School{
				ID: uint(schoolID),
			},
		},
	}

	studentID := c.QueryInt("studentId", 0)

	if roleID != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User can't access this page",
		})
	}

	dashboardData, err := dashboardController.dashboardService.GetDashboardFromParent(user, studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(dashboardData)
}

// @Summary Get Dashboard Data for Admin
// @Description Retrieve dashboard data for an admin user
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Success 200 {object} response.DashboardAdminResponse "Dashboard data for admin"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/dashboard/admin [get]
func (dashboardController *DashboardController) GetDashboardFromAdmin(c *fiber.Ctx) error {
	// Check if user has admin access
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	dashboardData, err := dashboardController.dashboardService.GetDashboardFromAdmin(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(dashboardData)
}
