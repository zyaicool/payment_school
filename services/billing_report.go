package services

import (
	"bytes"
	"fmt"
	"log"
	response "schoolPayment/dtos/response"
	"schoolPayment/repositories"
	"schoolPayment/utilities"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BillingReportService interface {
	GetBillingReport(page int, limit int, sortBy string, sortOrder string, schoolGradeId int, schoolClassId int, paymentStatusId int, schoolYearId int, bankAccountId int, billingType string, studentId int, userID int) (response.BillingReportResponse, error)
	ExportToExcel(limit int, page int, sortBy string, sortOrder string, schoolGradeId, schoolClassId, paymentStatus, schoolYearId, bankAccountId int, studentId int, billingType string, ctx *fiber.Ctx) (*bytes.Buffer, error)
}

type billingReportService struct {
	repo           repositories.BillingReportRepository
	userRepo       repositories.UserRepository
	userRepository repositories.UserRepository
}

func NewBillingReportService(repo repositories.BillingReportRepository, userRepo repositories.UserRepository) BillingReportService {
	return &billingReportService{repo: repo, userRepo: userRepo}
}

func (s *billingReportService) GetBillingReport(page int, limit int, sortBy string, sortOrder string, schoolGradeId int, schoolClassId int, paymentStatusId int, schoolYearId int, bankAccountId int, billingType string, studentId int, userID int) (response.BillingReportResponse, error) {
	var rsp response.BillingReportResponse
	var reportList []response.BillingReportDetail
	listBilling := response.ListReportBilling{
		Page:      page,
		Limit:     limit,
		TotalPage: 0,
		TotalData: 0,
		Data:      []response.BillingReportDetail{},
	}
	rsp.ListBillingReport = listBilling

	user, err := s.userRepo.GetUserByID(uint(userID))
	if err != nil {
		return rsp, err
	}

	summary, err := s.repo.GetSummaryBillingReport(schoolGradeId, schoolClassId, paymentStatusId, schoolYearId, bankAccountId, billingType, studentId, user)
	if err != nil {
		fmt.Println("Error in repository:", err)
		return rsp, err
	}

	rsp.TotalBillingAmount = utilities.FormatBigInt(summary.TotalBillingAmount)
	rsp.TotalPayAmount = utilities.FormatBigInt(summary.TotalPayAmount)
	rsp.TotalPayAmount = utilities.FormatBigInt(summary.TotalPayAmount)
	rsp.TotalNotPayAmount = utilities.FormatBigInt(summary.TotalNotPayAmount)
	rsp.TotalBillingPay = summary.TotalBillingPay
	rsp.TotalBillingNotPay = summary.TotalBillingNotPay
	rsp.TotalStudent = summary.TotalStudent

	if sortBy != "" {
		sortBy = utilities.ChangeStringSortByBillingReport(sortBy)
	}
	billingReports, totalPages, totalData, err := s.repo.GetBillingReport(page, limit, sortBy, sortOrder, schoolGradeId, schoolClassId, paymentStatusId, schoolYearId, bankAccountId, billingType, studentId, user)
	if err != nil {
		fmt.Println("Error in repository:", err)
		return rsp, err
	}

	listBilling.TotalPage = totalPages
	listBilling.TotalData = totalData

	for _, report := range billingReports {
		bankName, _ := GetBankName(report.BankName)
		bankAccountName := fmt.Sprintf("%s-%s", bankName, report.AccountNumber)

		responseBillingReport := response.BillingReportDetail{
			DetailBillingName: report.DetailBillingName,
			BillingType:       report.BillingType,
			StudentName:       report.StudentName,
			SchoolGradeName:   report.SchoolGradeName,
			SchoolClassName:   report.SchoolClassName,
			SchoolYearName:    report.SchoolYearName,
			Amount:            report.Amount,
			BankAccountName:   bankAccountName,
			PaymentStatus:     report.PaymentStatus,
		}

		reportList = append(reportList, responseBillingReport)
	}
	listBilling.Data = reportList
	rsp.ListBillingReport = listBilling
	return rsp, nil
}

func (s *billingReportService) ExportToExcel(limit int, page int, sortBy string, sortOrder string, schoolGradeId, schoolClassId, paymentStatus, schoolYearId, bankAccountId int, studentId int, billingType string, ctx *fiber.Ctx) (*bytes.Buffer, error) {
	userClaims := ctx.Locals("user").(jwt.MapClaims)
	userLoginID := int(userClaims["user_id"].(float64))

	// Ambil data rekapitulasi dan laporan billing
	report, err := s.GetBillingReport(page, limit, sortBy, sortOrder, schoolGradeId, schoolClassId, paymentStatus, schoolYearId, bankAccountId, billingType, studentId, userLoginID)

	if err != nil {
		log.Printf("Error getting payment report: %v", err)
		return nil, fmt.Errorf("failed to get payment report data: %v", err)
	}

	// Log jumlah data untuk Excel
	log.Printf("Successfully retrieved %d records for Excel export", len(report.ListBillingReport.Data))

	// Buat instance Excel utility
	excelUtil := utilities.NewExcelUtility()
	defer excelUtil.Close()

	sheetName := "BillingReport"
	excelUtil.File.SetSheetName("Sheet1", sheetName)

	// Buat gaya untuk teks tebal
	boldStyle, err := excelUtil.CreateBoldStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to create bold style: %v", err)
	}

	// Tulis header rekapitulasi di atas tabel
	fmt.Printf("Rp %s\n", report.TotalBillingAmount)
	fmt.Printf("Rp %s\n", utilities.FormatCurrency(report.TotalBillingAmount))
	excelUtil.SetCellValue(sheetName, "B1", "Total Tagihan")
	excelUtil.SetCellValue(sheetName, "C1", fmt.Sprintf(utilities.FormatCurrency(report.TotalBillingAmount)))
	excelUtil.SetCellValue(sheetName, "B2", "Total Sudah Dibayar")
	excelUtil.SetCellValue(sheetName, "C2", fmt.Sprintf(utilities.FormatCurrency(report.TotalPayAmount)))
	excelUtil.SetCellValue(sheetName, "B3", "Total Belum Dibayar")
	excelUtil.SetCellValue(sheetName, "C3", fmt.Sprintf(utilities.FormatCurrency(report.TotalNotPayAmount)))
	excelUtil.SetCellValue(sheetName, "B4", "Jumlah Sudah Dibayar")
	excelUtil.SetCellValue(sheetName, "C4", utilities.FormatWithThousandsSeparator(report.TotalBillingPay))
	excelUtil.SetCellValue(sheetName, "B5", "Jumlah Belum Dibayar")
	excelUtil.SetCellValue(sheetName, "C5", utilities.FormatWithThousandsSeparator(report.TotalBillingNotPay))
	excelUtil.SetCellValue(sheetName, "B6", "Jumlah Siswa")
	excelUtil.SetCellValue(sheetName, "C6", utilities.FormatWithThousandsSeparator(report.TotalStudent))

	// Terapkan gaya bold pada semua header
	excelUtil.SetCellStyle(sheetName, "B1", "C6", boldStyle)

	// Tulis header tabel
	headers := []string{
		"No.", "Nama Tagihan", "Tipe Tagihan", "Nama Siswa", "Unit",
		"Kelas", "Tahun Ajaran", "Jumlah Tagihan", "Rekening Bank",
		"Status",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s8", string(rune(65+i))) // Baris 8 untuk header tabel
		excelUtil.SetCellValue(sheetName, cell, header)
		excelUtil.SetCellStyle(sheetName, cell, cell, boldStyle)
	}

	// Tulis data billing report ke dalam tabel (mulai dari baris ke-9)
	for idx, detail := range report.ListBillingReport.Data {
		rowIdx := idx + 9 // Start from row 9
		data := []interface{}{
			idx + 1,
			detail.DetailBillingName,
			detail.BillingType,
			detail.StudentName,
			detail.SchoolGradeName,
			detail.SchoolClassName,
			detail.SchoolYearName,
			"RP " + utilities.FormatWithThousandsSeparator(int(detail.Amount)),
			detail.BankAccountName,
			detail.PaymentStatus,
		}
		for colIdx, value := range data {
			cell := fmt.Sprintf("%s%d", string(rune(65+colIdx)), rowIdx)
			excelUtil.SetCellValue(sheetName, cell, value)
		}
	}

	// Sesuaikan lebar kolom secara otomatis
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, col := range columns {
		if err := excelUtil.AutoFitColumn(sheetName, col); err != nil {
			return nil, fmt.Errorf("failed to auto-fit column %s: %v", col, err)
		}
	}

	// Simpan ke buffer
	buffer := new(bytes.Buffer)
	if err := excelUtil.Write(buffer); err != nil {
		return nil, fmt.Errorf("failed to write Excel file: %v", err)
	}

	return buffer, nil
}
