package controllers

import (
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	services "schoolPayment/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type LoginController struct {
	loginService services.LoginService
}

func NewLoginController(loginService services.LoginService) *LoginController {
	return &LoginController{loginService: loginService}
}

// @Summary Login
// @Description Authenticate a user and retrieve a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param firebaseToken header string false "Firebase token, only for mobile apps"
// @Param schoolId header string false "School ID, defaults to 0 if not provided"
// @Param loginRequest body request.LoginRequest true "Login request body"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid input or request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/login [post]
func (loginController *LoginController) Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest

	// Step 1: Parse the JSON request
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	firebaseToken := c.Get("firebaseToken")
	schoolIdStr := c.Get("schoolId", "0") // Set default value to "0" if schoolId is not provided
	schoolId, err := strconv.Atoi(schoolIdStr)
	if err != nil {
		// If conversion fails, set schoolId to 0
		schoolId = 0
	}

	// Step 2: Call the service to authenticate the user
	token, err := loginController.loginService.LoginService(loginRequest.Email, loginRequest.Password, firebaseToken, uint(schoolId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Step 3: Respond with the JWT token
	return c.JSON(response.LoginResponse{Token: token})
}

// @Summary Auth User
// @Description Get Data User After Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Success 200 {object} response.DataUserForAuth
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth [get]
func (loginController *LoginController) AuthUser(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	user, err := loginController.loginService.GetUserFromAuth(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

// @Summary Invalidate Firebase Token
// @Description Invalidate a Firebase token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param firebaseToken body string true "Firebase token to invalidate"
// @Param userId body int true "User ID"
// @Success 200 {object} response.InvalidateFirebaseTokenResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invalidFirebaseToken [post]
func (loginController *LoginController) InvalidateFirebaseToken(c *fiber.Ctx) error {
    var request struct {
        FirebaseToken string `json:"firebaseToken"`
        UserId        int    `json:"userId"`
    }

    // Parse request body
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
            Error: "Invalid request body",
        })
    }

    // Call the service method to invalidate the token
    err := loginController.loginService.UpdateFirebaseTokenStatus(request.FirebaseToken, request.UserId)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
            Error: err.Error(),
        })
    }

    // Return success response if the token is successfully invalidated
    return c.JSON(response.InvalidateFirebaseTokenResponse{
        Message: "Firebase token invalidated successfully",
    })
}

