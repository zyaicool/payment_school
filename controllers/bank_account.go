package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"
	services "schoolPayment/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type BankAccountController struct {
	bankAccountService services.BankAccountServiceInterface
}

func NewBankAccountController(bankAccountService services.BankAccountServiceInterface) BankAccountController {
	return BankAccountController{bankAccountService: bankAccountService}
}

// @Summary Create Bank Account
// @Description Create a new bank account for the authenticated user
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param bankAccount body request.BankAccountCreateRequest true "Bank Account Data"
// @Success 200 {object} map[string]interface{} "Bank account created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 500 {object} map[string]interface{} "Failed to create bank account"
// @Router /api/v1/bankAccount/create [post]
func (bankAccountController *BankAccountController) CreateBankAccount(c *fiber.Ctx) error {
	var bankRequest *request.BankAccountCreateRequest
	var userID uint = 0

	// Extract user ID from token claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = uint(userClaimID)
		}
	}

	if err := c.BodyParser(&bankRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createBankAccount, err := bankAccountController.bankAccountService.CreateBankAccount(bankRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data is saved successfully.",
		"data":    createBankAccount,
	})
}

// @Summary Update Bank Account
// @Description Update an existing bank account
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "Bank Account ID"
// @Param bankAccount body request.BankAccountCreateRequest true "Updated Bank Account Data"
// @Success 200 {object} map[string]interface{} "Bank account updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 409 {object} map[string]interface{} "Account number already exists"
// @Failure 500 {object} map[string]interface{} "Failed to update bank account"
// @Router /api/v1/bankAccount/update/{id} [put]
func (bankAccountController *BankAccountController) UpdateBankAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	bankID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bank account ID",
		})
	}

	var bankRequest request.BankAccountCreateRequest
	var userID int = 0

	fmt.Println("bankAccount", bankRequest)

	// Extract user ID from token claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if ok {
		if userClaimID, ok := userClaims["user_id"].(float64); ok {
			userID = int(userClaimID)
		}
	}

	if err := c.BodyParser(&bankRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedBankAccount, err := bankAccountController.bankAccountService.UpdateBankAccount(uint(bankID), &bankRequest, userID)
	if err != nil {
		if err.Error() == "account number already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data updated successfully.",
		"data":    updatedBankAccount,
	})
}

// @Summary Get All Bank Accounts
// @Description Retrieve all bank accounts for the authenticated user
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 0)"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc or desc)"
// @Success 200 {array} models.BankAccount "List of bank accounts"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve bank accounts"
// @Router /api/v1/bankAccount/getListBankAccount [get]
func (bankAccountController *BankAccountController) GetAllBankAccounts(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := uint(userClaims["user_id"].(float64))
	roleId := uint(userClaims["role_id"].(float64))
	schoolId := uint(userClaims["school_id"].(float64))

	var user = models.User{
		Master: models.Master{
			ID: userID,
		},
		RoleID: roleId,
		UserSchool: &models.UserSchool{
			School: &models.School{
				ID: schoolId,
			},
		},
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 0)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")

	listBankAccounts, err := bankAccountController.bankAccountService.GetAllBankAccounts(page, limit, search, sortBy, sortOrder, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(listBankAccounts)
}

// @Summary Get Bank Account by ID
// @Description Retrieve details of a specific bank account by its ID
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "Bank Account ID"
// @Success 200 {object} models.BankAccount "Bank account details"
// @Failure 400 {object} map[string]interface{} "Invalid bank account ID"
// @Failure 404 {object} map[string]interface{} "Bank account not found"
// @Failure 500 {object} map[string]interface{} "Failed to retrieve bank account"
// @Router /api/v1/bankAccount/detail/{id} [get]
func (bankAccountController *BankAccountController) GetBankAccountByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bank account ID",
		})
	}

	bankAccount, err := bankAccountController.bankAccountService.GetBankAccountDetails(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Bank account not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bank account",
		})
	}

	return c.JSON(bankAccount)
}

// @Summary Delete Bank Account
// @Description Delete a specific bank account by its ID
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "Bank Account ID"
// @Success 200 {object} map[string]interface{} "Bank account deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid bank account ID"
// @Failure 404 {object} map[string]interface{} "Bank account not found"
// @Failure 500 {object} map[string]interface{} "Failed to delete bank account"
// @Router /api/v1/bankAccount/delete/{id} [delete]
func (bankAccountController *BankAccountController) DeleteBankAccount(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	err = bankAccountController.bankAccountService.DeleteBankAccount(uint(id))
	if err != nil {
		if err.Error() == "bank account not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Bank account deleted successfully."})
}

// @Summary Get Bank Names
// @Description Retrieve a list of banks from a predefined JSON file
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Success 200 {object} response.BankResponse "List of bank names"
// @Failure 500 {object} map[string]interface{} "Failed to load bank names"
// @Router /api/v1/bankAccount/listBankName [get]
func (bankAccountController *BankAccountController) GetBankName(c *fiber.Ctx) error {
	banks, err := models.LoadBanks("data/banks.json") // Path to your JSON file
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load bank accounts",
		})
	}

	// Wrap the bank data in the desired structure
	var responseData []response.BankData
	for _, bank := range banks {
		responseData = append(responseData, response.BankData{
			CodeBank: bank.Code,
			Alias:    bank.Alias,
			BankName: bank.Name,
		})
	}

	response := response.BankResponse{Data: responseData}

	return c.JSON(response)
}
