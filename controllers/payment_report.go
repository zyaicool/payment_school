package controllers

import (
	services "schoolPayment/services"
	"schoolPayment/utilities"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PaymentReportController struct {
	paymentReportService services.PaymentReportService
}

func NewPaymentReportController(paymentReportService services.PaymentReportService) *PaymentReportController {
	return &PaymentReportController{
		paymentReportService: paymentReportService,
	}
}

// @Summary Get Payment Report
// @Description Retrieve a paginated list of payment reports based on various filters
// @Tags Payment Report
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param sortBy query string false "Field to sort by"
// @Param sortOrder query string false "Sort order (asc/desc)" default(asc)
// @Param paymentTypeId query int false "Filter by payment type ID"
// @Param userId query int false "Filter by user ID"
// @Param studentId query int false "Filter by student ID"
// @Param paymentStatus query string false "Filter by payment status"
// @Param schoolGradeId query int false "Filter by school grade ID"
// @Param schoolClassId query int false "Filter by school class ID"
// @Param startDate query string false "Start date filter in yyyy-MM-dd format"
// @Param endDate query string false "End date filter in yyyy-MM-dd format"
// @Success 200 {object} map[string]interface{} "Paginated payment report data"
// @Failure 400 {object} map[string]interface{} "Invalid input or date range"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/paymentReport/getList [get]
func (c *PaymentReportController) GetPaymentReport(ctx *fiber.Ctx) error {
	// Check access first
	err := utilities.CheckAccessUserTuKasirAdminSekolah(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := ctx.Locals("user").(jwt.MapClaims)
	userLoginID := int(userClaims["user_id"].(float64))
	schoolLoginId := int(userClaims["school_id"].(float64))

	// Get query parameters with defaults
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	sortBy := ctx.Query("sortBy", "")
	sortOrder := ctx.Query("sortOrder", "asc")
	paymentTypeId := ctx.QueryInt("paymentTypeId", 0)
	userId := ctx.QueryInt("userId", 0)
	studentId := ctx.QueryInt("studentId", 0)
	paymentStatus := ctx.Query("paymentStatus", "")
	schoolGradeId := ctx.QueryInt("schoolGradeId", 0)
	schoolClassId := ctx.QueryInt("schoolClassId", 0)

	var startDate, endDate time.Time

	// Parse startDate if provided
	startDateStr := ctx.Query("startDate", "")
	if startDateStr != "" {
		parsedStartDate, err := utilities.ParseDate(startDateStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format",
			})
		}
		startDate = parsedStartDate
	}

	// Parse endDate if provided
	endDateStr := ctx.Query("endDate", "")
	if endDateStr != "" {
		parsedEndDate, err := utilities.ParseDate(endDateStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end date format",
			})
		}
		endDate = parsedEndDate

		// Validate date range only if both dates are provided
		if !startDate.IsZero() && endDate.Before(startDate) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "End date must be after start date",
			})
		}
	}

	// Get report data
	report, err := c.paymentReportService.GetPaymentReport(
		page, limit, sortBy, sortOrder,
		paymentTypeId, userId, startDate, endDate,
		studentId, false, userLoginID, paymentStatus,
		schoolGradeId, schoolClassId, schoolLoginId,
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(report)
}

// @Summary Export Payment Report to Excel
// @Description Generate an Excel file of payment reports based on various filters
// @Tags Payment Report
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param Authorization header string true "Authorization" format("Bearer token")
// @Param paymentTypeId query int false "Filter by payment type ID"
// @Param userId query int false "Filter by user ID"
// @Param studentId query int false "Filter by student ID"
// @Param paymentStatus query string false "Filter by payment status"
// @Param schoolGradeId query int false "Filter by school grade ID"
// @Param schoolClassId query int false "Filter by school class ID"
// @Param startDate query string false "Start date filter in yyyy-MM-dd format"
// @Param endDate query string false "End date filter in yyyy-MM-dd format"
// @Success 200 {file} file "Excel file with payment report data"
// @Failure 400 {object} map[string]interface{} "Invalid input or date range"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/paymentReport/exportExcel [get]
func (c *PaymentReportController) ExportPaymentReportToExcel(ctx *fiber.Ctx) error {

	// Check access first
	err := utilities.CheckAccessUserTuKasirAdminSekolah(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get query parameters
	paymentTypeId := ctx.QueryInt("paymentTypeId", 0)
	userId := ctx.QueryInt("userId", 0)
	studentId := ctx.QueryInt("studentId", 0)
	paymentStatus := ctx.Query("paymentStatus", "")
	schoolGradeId := ctx.QueryInt("schoolGradeId", 0)
	schoolClassId := ctx.QueryInt("schoolClassId", 0)

	var startDate, endDate time.Time

	// Parse startDate if provided
	startDateStr := ctx.Query("startDate", "")
	if startDateStr != "" {
		parsedStartDate, err := utilities.ParseDate(startDateStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format",
			})
		}
		startDate = parsedStartDate
	}

	// Parse endDate if provided
	endDateStr := ctx.Query("endDate", "")
	if endDateStr != "" {
		parsedEndDate, err := utilities.ParseDate(endDateStr)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end date format",
			})
		}
		endDate = parsedEndDate

		// Validate date range
		if !startDate.IsZero() && endDate.Before(startDate) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "End date must be after start date",
			})
		}
	}

	// Generate Excel file
	buffer, err := c.paymentReportService.ExportToExcel(paymentTypeId, userId, startDate, endDate, studentId, ctx, paymentStatus, schoolGradeId, schoolClassId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Send(buffer.Bytes())
}
