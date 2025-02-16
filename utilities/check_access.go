package utilities

import (
	"fmt"
	"schoolPayment/constants"
	"schoolPayment/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAccessAdminSekolah(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := repositories.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID != 5 {
		return fmt.Errorf(constants.MessageUserCantAccessPage)
	}

	return nil
}

func CheckAccessToBillingType(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := repositories.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID == 1 || roleID == 5 {
		return nil
	}

	return fmt.Errorf(constants.MessageUserCantAccessPage)
}

func CheckAccessUserOrtu(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := repositories.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID != 2 {
		return fmt.Errorf(constants.MessageUserCantAccessPage)
	}

	return nil
}

func CheckAccessUserTuKasirAdminSekolah(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))

	// Allow access for Kasir (4), and Tata Usaha (3), Admin Sekolah (5)
	if roleID == 3 || roleID == 4 || roleID == 5 {
		return nil
	}

	return fmt.Errorf(constants.MessageUserCantAccessPage)
}

func CheckAccessExceptUserOrtu(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	roleID := int(userClaims["role_id"].(float64))
	_, err := repositories.GetRoleByID(uint(roleID))
	if err != nil {
		return err
	}

	if roleID == 2 {
		return fmt.Errorf(constants.MessageUserCantAccessPage)
	}

	return nil
}
