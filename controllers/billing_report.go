package controllers

import (
	"fmt"
	"log"
	"schoolPayment/services"
	"schoolPayment/utilities"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BillingReportController struct {
	service services.BillingReportService
}

func NewBillingReportController(service services.BillingReportService) *BillingReportController {
	return &BillingReportController{service: service}
}

// @Summary Get Billing Reports
// @Description Retrieve a list of billing reports based on various filters and pagination
// @Tags Billing Reports
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of records per page (default: 10)"
// @Param sortBy query string false "Field to sort by (e.g., created_at)"
// @Param sortOrder query string false "Sort order (asc or desc, default: ASC)"
// @Param schoolGradeId query int false "School grade ID filter"
// @Param schoolClassId query int false "School class ID filter"
// @Param paymentStatusId query int false "Payment status ID filter"
// @Param schoolYearId query int false "School year ID filter"
// @Param bankAccountId query int false "Bank account ID filter"
// @Param billingType query string false "Billing type filter"
// @Param studentId query int false "Student ID filter"
// @Success 200 {object} map[string]interface{} "Paginated billing report data"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to get billing reports"
// @Router /api/v1/billingReport/getList [get]
func (c *BillingReportController) GetBillingReports(ctx *fiber.Ctx) error {
	err := utilities.CheckAccessExceptUserOrtu(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userClaims := ctx.Locals("user").(jwt.MapClaims)
	userID := int(userClaims["user_id"].(float64))

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit number"})
	}
	sortBy := ctx.Query("sortBy", "")
	sortOrder := ctx.Query("sortOrder", "ASC")
	schoolGradeId, _ := strconv.Atoi(ctx.Query("schoolGradeId", "0"))
	schoolClassId, _ := strconv.Atoi(ctx.Query("schoolClassId", "0"))
	paymentStatusId, _ := strconv.Atoi(ctx.Query("paymentStatusId", "0"))
	schoolYearId, _ := strconv.Atoi(ctx.Query("schoolYearId", "0"))
	bankAccountId, _ := strconv.Atoi(ctx.Query("bankAccountId", "0"))
	billingType := ctx.Query("billingType", "")
	studentId, _ := strconv.Atoi(ctx.Query("studentId", "0"))

	response, err := c.service.GetBillingReport(page, limit, sortBy, sortOrder, schoolGradeId, schoolClassId, paymentStatusId, schoolYearId, bankAccountId, billingType, studentId, userID)
	if err != nil {
		fmt.Println("Error in service:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get billing report"})
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary Export Billing Report to Excel
// @Description Export billing reports to an Excel file based on specified filters
// @Tags Billing Reports
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param schoolGradeId query int false "School grade ID filter"
// @Param schoolClassId query int false "School class ID filter"
// @Param paymentStatusId query int false "Payment status ID filter"
// @Param schoolYearId query int false "School year ID filter"
// @Param bankAccountId query int false "Bank account ID filter"
// @Param billingType query string false "Billing type filter"
// @Param studentId query int false "Student ID filter"
// @Success 200 {file} file "Excel file containing the exported billing report"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to export billing report"
// @Router /api/v1/billingReport/exportExcel [get]
func (c *BillingReportController) ExportBillingReportToExcel(ctx *fiber.Ctx) error {
	// Ambil user dari context
	err := utilities.CheckAccessExceptUserOrtu(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims, ok := ctx.Locals("user").(jwt.MapClaims)
	if !ok {
		log.Println("Unauthorized: User claims not found")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: User not found"})
	}

	// Validasi user ID
	userIDFloat, ok := userClaims["user_id"].(float64)
	if !ok {
		log.Println("Unauthorized: Invalid user ID format")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid user ID"})
	}
	userID := int(userIDFloat)
	log.Printf("Export request by user ID: %d", userID)

	// Ambil parameter query
	schoolGradeId, _ := strconv.Atoi(ctx.Query("schoolGradeId", "0"))
	schoolClassId, _ := strconv.Atoi(ctx.Query("schoolClassId", "0"))
	paymentStatusId, _ := strconv.Atoi(ctx.Query("paymentStatusId", "0"))
	schoolYearId, _ := strconv.Atoi(ctx.Query("schoolYearId", "0"))
	bankAccountId, _ := strconv.Atoi(ctx.Query("bankAccountId", "0"))
	billingType := ctx.Query("billingType", "")
	studentId, _ := strconv.Atoi(ctx.Query("studentId", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	sortBy := ctx.Query("sortBy", "")
	sortOrder := ctx.Query("sortOrder", "ASC")

	// Panggil service untuk ekspor
	buffer, err := c.service.ExportToExcel(limit, page, sortBy, sortOrder, schoolGradeId, schoolClassId, paymentStatusId, schoolYearId, bankAccountId, studentId, billingType, ctx)
	if err != nil {
		log.Printf("Error in service ExportToExcel: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export billing report"})
	}

	// limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	// if err != nil || limit <= 0 {
	// 	limit = 10 // Nilai default jika parameter tidak valid
	// }

	// Set header untuk file Excel
	currentTime := time.Now()
	formattedDate := currentTime.Format("02-01-2006")
	formattedTime := currentTime.Format("15.04")
	filename := fmt.Sprintf("Laporan Tagihan %s %s.xlsx", formattedDate, formattedTime)
	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Kirim buffer sebagai response
	return ctx.SendStream(buffer)
}
