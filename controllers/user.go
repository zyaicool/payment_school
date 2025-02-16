package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// @Summary Get All Users
// @Description Get a list of users with optional filters for role, pagination, and status
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param roleId query string false "Comma-separated list of role IDs"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search term"
// @Param schoolId query int false "School ID"
// @Param sortBy query string false "Sort field" default("created_at")
// @Param sortOrder query string false "Sort order" default("desc")
// @Param status query string false "Status filter" enum(true,false)
// @Success 200 {array} response.UserListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/getAllUser [get]
func (userController *UserController) GetAllUser(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	roleIDParam := c.Query("roleId")
	var roleIDs []int

	if roleIDParam != "" {
		roleIDString := strings.Split(roleIDParam, ",")
		for _, roleIDStr := range roleIDString {
			roleID, err := strconv.Atoi(roleIDStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid roleId parameter",
				})
			}

			roleIDs = append(roleIDs, roleID)
		}
	}

	// get query params pagination and email
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")

	schoolID := c.QueryInt("schoolId", 0)
	sortBy := c.Query("sortBy", "created_at")
	sortOrder := c.Query("sortOrder", "desc")
	var status *bool
	statusParam := c.Query("status")
	if statusParam != "" {
		parsedStatus, err := strconv.ParseBool(statusParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status parameter, must be true or false",
			})
		}
		status = &parsedStatus
	}

	users, err := userController.userService.GetAllUser(page, limit, search, roleIDs, userID, schoolID, sortBy, sortOrder, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(users)
}

// @Summary Get User Data By ID
// @Description Get detailed information about a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "User ID"
// @Success 200 {object} response.DetailUserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/detail/{id} [get]
func (userController *UserController) GetDataUser(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	//user, err := userController.userService.GetUserByID(uint(id))
	user, err := userController.userService.GetUserDetail(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}

	return c.JSON(user)
}

// @Summary Get User Data By Email
// @Description Get detailed information about a user by their email address
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "Email address of the user"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/getByEmail [get]
func GetDataByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	user, err := services.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(user)
}

// @Summary Create a new user
// @Description Create a new user after checking access and parsing the request body
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param user body request.UserCreateRequest true "User Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/create [post]
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user *request.UserCreateRequest
	var userID int = 0

	// Get user information from the token claims stored in context
	userClaims, ok := c.Locals("user").(jwt.MapClaims)

	if ok {
		// If the user claims exist and are of type jwt.MapClaims, extract the user ID.
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createUser, err := uc.userService.CreateUser(user, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createUser,
	})
}

// @Summary Update user data
// @Description Update an existing user with the given data
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "User ID"
// @Param user body request.UserUpdateRequest true "Updated User Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/edit/{id} [put]
func (userController *UserController) UpdateUser(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	id, _ := c.ParamsInt("id")
	var user *request.UserUpdateRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	// Panggil service untuk update password
	updatedUser, err := userController.userService.UpdateUserService(uint(id), user, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedUser,
	})
}

// @Summary Delete a user
// @Description Delete a user with the specified ID
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/delete/{id} [delete]
func (userController *UserController) DeleteUser(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	user, err := userController.userService.GetUserByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}

	_, err = userController.userService.DeleteUserService(&user, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}

// @Summary Verify email using token
// @Description Verify the user's email address by using a token
// @Tags User
// @Accept json
// @Produce json
// @Param token query string true "Token for email verification"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/verifyEmail [get]
func (userController *UserController) EmailVerification(c *fiber.Ctx) error {
	token := c.Query("token")
	err := userController.userService.ValidateEmail(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email has verification.",
	})
}

// @Summary Resend Email Verification
// @Description Resend Email Verification to user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/resendVerificationEmail/{id} [post]
func (userController *UserController) ResendEmailVerification(c *fiber.Ctx) error {

	userId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	err = userController.userService.ResendEmailVerification(userId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email send.",
	})
}

// @Summary Generate Token for Change Password
// @Description Generate a token for changing password via email
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/generateTokenChangePassword [get]
func (userController *UserController) GenerateTokenChangePassword(c *fiber.Ctx) error {
	email := c.Query("email")
	err := userController.userService.GenerateTokenChangePassword(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Change password link has sended to your email, please imidiately checking that.",
	})
}

// @Summary Verify Token for Change Password
// @Description Verify the token for change password
// @Tags Users
// @Accept json
// @Produce json
// @Param token query string true "Change Password Token"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/verifyTokenChangePassword [get]
func (userController *UserController) VerifyTokenChangePassword(c *fiber.Ctx) error {
	token := c.Query("token")
	err := userController.userService.VerifyTokenChangePassword(token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		} else {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token is valid.",
	})
}

// @Summary Change Password
// @Description Change user password using the provided token
// @Tags Users
// @Accept json
// @Produce json
// @Param token query string true "Token"
// @Param user body request.ChangePasswordRequest true "New Password"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/changePassword [post]
func (userController *UserController) ChangePassword(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	token := c.Query("token")
	var request request.ChangePasswordRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	// Panggil service untuk update password
	_, err := userController.userService.ChangePassword(token, request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Change password success.",
	})
}

// @Summary Download Excel Format for User
// @Description Generate and download an Excel file for a user
// @Tags Users
// @Accept json
// @Produce octet-stream
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Success 200 {file} file "Excel file"
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/getFileExcel [get]
func (userController *UserController) DownloadFileExcelFormatForUser(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	fmt.Println("userId", userID)

	buffer, err := userController.userService.GenerateFileExcelForUser(c, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStream(buffer)
}

// @Summary Upload Users via File
// @Description Upload users via CSV or Excel file
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param file formData file true "File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/upload [post]
func (userController *UserController) UploadUsers(c *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	// Check file type (CSV or Excel)
	extension := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	switch extension {
	case ".xlsx", ".xls":
		return userController.userService.HandleExcelFileUser(file, c, userID)
	case ".csv":
		return userController.userService.HandleCSVFileUser(file, c, userID)
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unsupported file type",
		})
	}
}

// @Summary Check if Email is Already Taken
// @Description Check if the provided email is available or already taken
// @Tags Users
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Success 200 {object} map[string]bool
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/checkEmaill [get]
func (userController *UserController) CheckExistingEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	valid := userController.userService.CheckEmail(email)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"available": valid,
	})
}

// @Summary Check if Username is Already Taken
// @Description Check if the provided username is available or already taken
// @Tags Users
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Success 200 {object} map[string]bool
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/checkUsername [get]
func CheckExistingUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	valid := services.CheckUsername(username)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"available": valid,
	})
}

// @Summary Upload User Photo
// @Description Upload a new profile photo for the user
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param image formData file true "Profile Image"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/uploadFoto [post]
func (userController *UserController) UploadUserPhoto(c *fiber.Ctx) error {
	// Get user claims from JWT token
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Upload the image and get the filename
	newFileName, err := utilities.UploadImage(file, c, "user_image", "image")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create the request object
	request := request.UpdateUserImageRequest{
		Image: newFileName,
	}

	// Call service to update user photo
	updatedUser, err := userController.userService.UpdateUserPhotoService(request, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Photo profile has been updated successfully",
		"data":    updatedUser,
	})
}

// @Summary Change User Password without Token
// @Description Change the user's password without using a token, using the current password and new password
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body request.UpdatePasswordRequest true "Password Change Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/user/updatePassword/{id} [put]
func (userController *UserController) ChangePasswordWithoutToken(c *fiber.Ctx) error {
	// Get user information from the token claims stored in context
	userIDStr := c.Params("id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Panggil service untuk update password
	var request request.UpdatePasswordRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	_, err = userController.userService.ChangePasswordWithoutToken(uint(userID), request.OldPassword, request.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Change password success.",
	})
}
