package controllers

import (
	"log"

	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	services "schoolPayment/services"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BillingController struct {
	billingService services.BillingServiceInterface
}

func NewBillingController(billingService services.BillingServiceInterface) *BillingController {
	return &BillingController{billingService: billingService}
}

// @Summary Get All Billings
// @Description Retrieve all billing data with optional filters
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param billingType query string false "Billing type"
// @Param paymentType query string false "Payment type"
// @Param schoolGrade query string false "School grade"
// @Param bankAccountId query int false "Bank account ID"
// @Param sort query string false "Sort order (asc/desc)"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Order of sorting (asc/desc)"
// @Param isDonation query bool false "Filter for donations"
// @Success 200 {object} response.BillingListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/getAllBilling [get]
func (billingController *BillingController) GetAllBilling(c *fiber.Ctx) error {
	// err := CheckAccessUserTuKasirAdminSekolah(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search")
	billingType := c.Query("billingType")
	paymentType := c.Query("paymentType")
	schoolGrade := c.Query("schoolGrade")
	bankAccountId := c.QueryInt("bankAccountId", 0)
	sort := c.Query("sort", "desc")
	sortBy := c.Query("sortBy", "") // Default sort field
	sortOrder := c.Query("sortOrder", "asc")
	isDonation := c.QueryBool("isDonation", false)

	listBilling, err := billingController.billingService.GetAllBilling(page, limit, search, billingType, paymentType, schoolGrade, sort, sortBy, sortOrder, bankAccountId, &isDonation, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch data",
		})
	}
	return c.JSON(listBilling)
}

// @Summary Get Billing By ID
// @Description Retrieve billing data by its ID
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Billing ID"
// @Success 200 {object} response.BillingDetailResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/billing/detail/{id} [get]
func (billingController *BillingController) GetDataBilling(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, _ := c.ParamsInt("id")
	billing, err := billingController.billingService.GetBillingByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(billing)
}

// @Summary Create a New Billing
// @Description Create a new billing entry
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param request body request.BillingCreateRequest true "Billing creation payload"
// @Success 200 {object} models.Billing
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/create [post]
func (billingController *BillingController) CreateBilling(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var billingRequest *request.BillingCreateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	if err := c.BodyParser(&billingRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	createBilling, err := billingController.billingService.CreateBilling(billingRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil disimpan.",
		"data":    createBilling,
	})
}

// @Summary Update Billing
// @Description Update an existing billing entry
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Billing ID"
// @Param request body request.BillingUpdateRequest true "Billing update payload"
// @Success 200 {object} models.Billing
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/update/{id} [put]
func UpdateBilling(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var billingRequest *request.BillingUpdateRequest

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	// Parse JSON dari request body ke DTO
	if err := c.BodyParser(&billingRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.CannotParseJsonMessage,
		})
	}

	updatedBilling, err := services.UpdateBilling(&services.BillingService{}, uint(id), billingRequest, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dirubah.",
		"data":    updatedBilling,
	})
}

// @Summary Delete Billing
// @Description Delete a billing entry by ID
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param id path int true "Billing ID"
// @Success 200 {object} models.Billing
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/delete/{id} [delete]
func DeleteBilling(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user information from the token claims stored in context
	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))
	id, _ := c.ParamsInt("id")

	_, err = services.DeleteBilling(&services.BillingService{}, uint(id), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data berhasil dihapus.",
	})
}

// @Summary Get All Billings
// @Description Retrieve all billing data with optional filters
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param billingType query string false "Billing type"
// @Param paymentType query string false "Payment type"
// @Param schoolGrade query string false "School grade"
// @Param bankAccountId query int false "Bank account ID"
// @Param sort query string false "Sort order (asc/desc)"
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Order of sorting (asc/desc)"
// @Param isDonation query bool false "Filter for donations"
// @Success 200 {object} []models.BillingStatus
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/billingStatus [get]
func (billingController *BillingController) GetBillingStatuses(c *fiber.Ctx) error {
	// Call the GetPaymentStatuses function in the service
	paymentStatuses, err := billingController.billingService.GetBillingStatuses("data/billing_status.json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch payment statuses",
		})
	}

	return c.JSON(fiber.Map{
		"data": paymentStatuses,
	})
}

func (controller *PaymentController) CheckBillingStatus(c *fiber.Ctx) error {
	transactionID := c.Params("transaction_id")
	// Logika untuk memeriksa status pembayaran dengan transactionID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Payment status for transaction " + transactionID,
	})
}

func GenerateInstallment(c *fiber.Ctx) error {
	err := utilities.CheckAccessUserTuKasirAdminSekolah(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, _ := c.ParamsInt("id")
	billing, err := services.GenerateInstallment(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": constants.DataNotFoundMessage,
		})
	}
	return c.JSON(billing)
}

// @Summary Create Donation
// @Description Create a donation billing entry
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param request body map[string]interface{} true "Donation creation payload"
// @Success 200 {object} response.BillingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/createDonation [post]
func (billingController *BillingController) CreateDonation(c *fiber.Ctx) error {

	var request struct {
		BillingName   string `json:"billingName"`
		SchoolGradeId int    `json:"schoolGradeId"`
		BankAccountId int    `json:"bankAccountId"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	userClaims := c.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	billingResponse, err := billingController.billingService.CreateDonation(request.BillingName, request.SchoolGradeId, request.BankAccountId, uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    billingResponse,
		"message": "Donation created successfully",
	})
}

// @Summary Get Billing By Student ID
// @Description Retrieve billing information for a specific student
// @Tags Billing
// @Accept json
// @Produce json
// @Param Authorization header string false "Authorization" format("Bearer token")
// @Param studentId query int true "Student ID"
// @Param schoolYearId query int false "School Year ID"
// @Param schoolGradeId query int false "School Grade ID"
// @Param schoolClassId query int false "School Class ID"
// @Success 200 {object} response.BillingByStudentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/billing/getBillingByStudentID [get]
func (c *BillingController) GetBillingByStudentID(ctx *fiber.Ctx) error {
	studentID := ctx.QueryInt("studentId", 0)
	schoolYearID := ctx.QueryInt("schoolYearId", 0)
	schoolGradeID := ctx.QueryInt("schoolGradeId", 0)
	schoolClassID := ctx.QueryInt("schoolClassId", 0)

	if studentID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "studentId are required",
		})
	}

	// Add validation for student existence if needed
	result, err := c.billingService.GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID)
	if err != nil {
		// Log the actual error for debugging
		log.Printf("Error getting billing by student ID: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve billing information",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
