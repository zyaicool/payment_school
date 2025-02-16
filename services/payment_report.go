package services

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"schoolPayment/constants"
	"schoolPayment/dtos/response"
	"schoolPayment/repositories"
	"schoolPayment/utilities"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PaymentReportService interface {
	GetPaymentReport(page, limit int, sortBy, sortOrder string, paymentTypeId, userId int, startDate, endDate time.Time, studentId int, isAllData bool, userLoginID int, paymentStatus string, schoolGradeId int, schoolClassId, schoolLoginId int) (response.PaymentReportResponse, error)
	ExportToExcel(paymentTypeId, userId int, startDate, endDate time.Time, studentId int, ctx *fiber.Ctx, paymentStatus string, schoolGradeId int, schoolClassId int) (*bytes.Buffer, error)
}

type paymentReportService struct {
	paymentReportRepository repositories.PaymentReportRepository
	userRepository          repositories.UserRepository
}

func NewPaymentReportService(paymentReportRepo repositories.PaymentReportRepository, userRepo repositories.UserRepository) PaymentReportService {
	return &paymentReportService{
		paymentReportRepository: paymentReportRepo,
		userRepository:          userRepo,
	}
}

func (s *paymentReportService) GetPaymentReport(
	page, limit int,
	sortBy, sortOrder string,
	paymentTypeId, userId int,
	startDate, endDate time.Time,
	studentId int, isAllData bool, userLoginID int, paymentStatus string, schoolGradeId int, schoolClassId, schoolLoginId int,
) (response.PaymentReportResponse, error) {
	var report response.PaymentReportResponse

	// Set default values if needed
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	if sortBy != "" {
		sortBy = utilities.ChangeStringSortByPaymentReport(sortBy)
		sortBy = utilities.ToSnakeCase(sortBy)
	}

	// Get detailed report data
	listReport, totalPage, totalData, err := s.paymentReportRepository.GetPaymentReport(
		page, limit, sortBy, sortOrder,
		paymentTypeId, userId, startDate, endDate,
		studentId, schoolLoginId, isAllData, paymentStatus,
		schoolGradeId, schoolClassId,
	)
	if err != nil {
		return report, fmt.Errorf("failed to get payment report: %v", err)
	}

	// Get summary data
	summary, err := s.paymentReportRepository.GetPaymentReportSummary(
		paymentTypeId, userId, startDate, endDate,
		studentId, schoolLoginId, paymentStatus,
		schoolGradeId, schoolClassId,
	)
	if err != nil {
		return report, fmt.Errorf("failed to get payment summary: %v", err)
	}

	// Construct response
	var totalTransactionAmount *big.Int
	if summary.TotalTransactionAmount == nil {
		totalTransactionAmount = big.NewInt(0) // Default to 0
	} else {
		totalTransactionAmount = summary.TotalTransactionAmount
	}
	report = response.PaymentReportResponse{
		TotalTransactionAmount: totalTransactionAmount,
		TotalTransaction:       summary.TotalTransaction,
		TotalStudent:           summary.TotalStudent,
		ListPaymentReport: response.ListPaymentReport{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: totalData,
			Data:      listReport,
		},
	}

	return report, nil
}

func (s *paymentReportService) ExportToExcel(paymentTypeId, userId int, startDate, endDate time.Time, studentId int, ctx *fiber.Ctx, paymentStatus string, schoolGradeId int, schoolClassId int) (*bytes.Buffer, error) {
	userClaims := ctx.Locals("user").(jwt.MapClaims)
	userLoginID := int(userClaims["user_id"].(float64))
	schoolLoginId := int(userClaims["school_id"].(float64))

	user, err := s.userRepository.GetUserByID(uint(userLoginID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %v", err)
	}

	report, err := s.GetPaymentReport(0, 0, "", "", paymentTypeId, userId, startDate, endDate, studentId, true, userLoginID, paymentStatus, schoolGradeId, schoolClassId, schoolLoginId)
	if err != nil {
		log.Printf("Error getting payment report: %v", err)
		return nil, fmt.Errorf("failed to get payment report data: %v", err)
	}

	// Log successful data retrieval
	log.Printf("Successfully retrieved %d records for Excel export", len(report.ListPaymentReport.Data))

	// Create Excel utility instance
	excelUtil := utilities.NewExcelUtility()
	defer excelUtil.Close()

	// Rename the default sheet to "Payment Report"
	excelUtil.File.SetSheetName("Sheet1", constants.PaymentReportParam)

	// Create center alignment style for "No." column
	centerStyle, err := excelUtil.CreateCenterStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create center style: %v", err)
	}

	// Create bold style for summary headers
	boldStyle, err := excelUtil.CreateBoldStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create bold style: %v", err)
	}

	// Create Rupiah style
	rupiahStyle, err := excelUtil.CreateRupiahStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create Rupiah style: %v", err)
	}

	// Write summary at the top
	summaryHeaders := []string{
		"Total Tagihan",
		"Jumlah Transaksi Sudah Dibayar",
		"Jumlah Siswa",
	}
	summaryValues := []string{
		utilities.RupiahFormat(report.TotalTransactionAmount), // Convert big.Int to string
		fmt.Sprintf("%d", report.TotalTransaction),
		fmt.Sprintf("%d", report.TotalStudent),
	}

	// Write summary headers and values
	for i, header := range summaryHeaders {
		row := i + 2
		headerCell := fmt.Sprintf("B%d", row)
		valueCell := fmt.Sprintf("C%d", row)

		excelUtil.SetCellValue(constants.PaymentReportParam, headerCell, header)
		excelUtil.SetCellStyle(constants.PaymentReportParam, headerCell, headerCell, boldStyle)

		excelUtil.SetCellValue(constants.PaymentReportParam, valueCell, summaryValues[i])
		excelUtil.SetCellStyle(constants.PaymentReportParam, valueCell, valueCell, rupiahStyle)
	}

	// Add empty rows between summary and table
	startTableRow := len(summaryHeaders) + 4

	// Define headers
	headers := []string{
		"No.", "Nomor Invoice", "Nama Siswa", "Tanggal Pembayaran", "Metode Pembayaran", "Nama Kasir",
		"Unit", "Kelas", "Total Pembayaran", "Status",
	}

	// Write headers with bold style starting after summary
	headerRow := startTableRow
	for col, header := range headers {
		cell := fmt.Sprintf("%s%d", utilities.GetColumnName(col), headerRow)
		excelUtil.SetCellValue(constants.PaymentReportParam, cell, header)
		excelUtil.SetCellStyle(constants.PaymentReportParam, cell, cell, boldStyle)
	}

	// Write data (update starting row to be after headers)
	for i, item := range report.ListPaymentReport.Data {
		row := headerRow + i + 1
		// Add index number starting from 1
		noCell := fmt.Sprintf("A%d", row)
		excelUtil.SetCellValue(constants.PaymentReportParam, noCell, i+1)
		excelUtil.SetCellStyle(constants.PaymentReportParam, noCell, noCell, centerStyle)
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("B%d", row), item.InvoiceNumber)
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("C%d", row), item.StudentName)
		if item.PaymentDate != nil {
			excelUtil.SetCellValue(
				constants.PaymentReportParam,
				fmt.Sprintf("D%d", row),
				item.PaymentDate.Format("02/01/2006 15:04:05"),
			)
		} else {
			excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("D%d", row), "")
		}
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("E%d", row), item.PaymentMethod)
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("F%d", row), strings.Title(item.Username)) //strings.Title(item.Username)
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("G%d", row), item.SchoolGradeName)
		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("H%d", row), item.SchoolClassName)

		// Set amount with Rupiah style
		amountCell := fmt.Sprintf("I%d", row)
		excelUtil.SetCellValue(constants.PaymentReportParam, amountCell, item.TotalAmount)
		excelUtil.SetCellStyle(constants.PaymentReportParam, amountCell, amountCell, rupiahStyle)

		excelUtil.SetCellValue(constants.PaymentReportParam, fmt.Sprintf("J%d", row), strings.Title(item.TransactionStatus))
	}

	// Adjust column widths automatically
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, col := range columns {
		if err := excelUtil.AutoFitColumn(constants.PaymentReportParam, col); err != nil {
			return nil, fmt.Errorf("failed to auto-fit column %s: %v", col, err)
		}
	}

	// // Set column widths
	// columnWidths := map[string]float64{
	// 	"A": 5,  // No.
	// 	"B": 30, // Detail Billing Name
	// 	"C": 20, // Billing Type
	// 	"D": 25, // Student Name
	// 	"E": 15, // School Grade
	// 	"F": 15, // School Class
	// 	"G": 15, // School Year
	// 	"H": 20, // Amount
	// 	"I": 20, // Bank Account
	// 	"J": 15, // Payment Status
	// 	"L": 25, // Summary column
	// 	"M": 20, // Summary values
	// }
	// if err := excelUtil.SetColumnWidths(constants.PaymentReportParam, columnWidths); err != nil {
	// 	return nil, fmt.Errorf("failed to set column widths: %v", err)
	// }

	// Update autofilter to start from header row
	if err := excelUtil.AddAutoFilter(constants.PaymentReportParam, fmt.Sprintf("A%d", headerRow), fmt.Sprintf("J%d", headerRow+len(report.ListPaymentReport.Data))); err != nil {
		return nil, fmt.Errorf("failed to add auto filter: %v", err)
	}

	// Write to buffer
	buffer, err := excelUtil.WriteToBuffer()
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
		return nil, fmt.Errorf("failed to write to buffer: %v", err)
	}
	filename := fmt.Sprintf("payment_report_%s.xlsx", user.UserSchool.School.SchoolName)
	ctx.Set(constants.ContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return buffer, nil
}
