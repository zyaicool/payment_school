package utilities

import (
	"fmt"
	"os"
	"strings"
	"time"

	request "schoolPayment/dtos/request"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// verify jwt token
func JWTProtected(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Please Login",
		})
	}

	// Extract the token part (after "Bearer ")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Save the claims to the context for later use
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Extract user data from claims
	// userID := (*claims)["UserID"]
	// roleID := (*claims)["RoleID"]
	c.Locals("user", *claims)
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	roleID := int(userClaims["role_id"].(float64))

	// Validate the user's role_id and account status
	userData, err := repositories.GetUserByID2(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if userData.RoleID != uint(roleID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized role access",
		})
	}
	if userData.IsBlock {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Account is blocked",
		})
	}

	// Proceed to the next handler
	return c.Next()
}

// generate jwt token
func GenerateJWT(user models.User, firebaseID string) (string, error) {
	// Create JWT claims with user information
	claims := &request.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
		SchoolID: func() uint {
			if user.UserSchool != nil {
				return user.UserSchool.SchoolID
			}
			return 0
		}(),
		FirebaseToken: firebaseID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token with the secret key and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("error signing the token: %v", err)
	}

	return tokenString, nil
}
