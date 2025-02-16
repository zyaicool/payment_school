package repositories

import (
	"fmt"
	"schoolPayment/configs"
	"schoolPayment/constants"
	"schoolPayment/dtos/response"
	"schoolPayment/models"
)

type BillingReportRepository interface {
	GetBillingReport(page int, limit int, sortBy string, sortOrder string, schoolGradeId int, schoolClassId int,
		paymentStatusId int, schoolYearId int, bankAccountId int, billingType string, studentId int, user models.User) ([]response.BillingReport, int, int64, error)
	GetSummaryBillingReport(schoolGradeId int, schoolClassId int, paymentStatusId int, schoolYearId int,
		bankAccountId int, billingType string, studentId int, user models.User) (response.BillingReportSummary, error)
}

type billingReportRepository struct{}

func NewBillingReportRepository() BillingReportRepository {
	return &billingReportRepository{}
}

func (r *billingReportRepository) GetBillingReport(page int, limit int, sortBy string, sortOrder string, schoolGradeId int, schoolClassId int, paymentStatusId int, schoolYearId int, bankAccountId int, billingType string, studentId int, user models.User) ([]response.BillingReport, int, int64, error) {
	var billingReports []response.BillingReport
	var totalData int64
	totalPages := 0
	offset := (page - 1) * limit

	// Validasi input pagination
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit jika tidak diberikan
	}

	// Query untuk mendapatkan total data tanpa LIMIT dan OFFSET
	countQuery := `
		SELECT COUNT(*)
		FROM billing_students bs
		LEFT JOIN billings b ON b.id = bs.billing_id
		LEFT JOIN students s ON s.id = bs.student_id
		LEFT JOIN school_grades sg ON sg.id = s.school_grade_id
		LEFT JOIN school_years sy ON sy.id = b.school_year_id  
		LEFT JOIN school_classes sc ON sc.id = s.school_class_id
		JOIN user_students us ON us.student_id = s.id
		JOIN user_schools usc ON usc.user_id = us.user_id
		WHERE bs.deleted_at IS NULL 
	`
	if user.UserSchool != nil {
		countQuery += fmt.Sprintf(constants.FilterByUscSchoolId, user.UserSchool.SchoolID)
	}
	if schoolClassId != 0 {
		countQuery += fmt.Sprintf(constants.FilterByUsSchoolClassId, schoolClassId)
	}
	if schoolGradeId != 0 {
		countQuery += fmt.Sprintf(constants.FilterBySgSchoolGradeId, schoolGradeId)
	}
	if paymentStatusId != 0 {
		countQuery += fmt.Sprintf(constants.FilterByBsPaymentStatus, paymentStatusId)
	}
	if schoolYearId != 0 {
		countQuery += fmt.Sprintf(constants.FilterBySySchoolYearId, schoolYearId)
	}
	if bankAccountId != 0 {
		countQuery += fmt.Sprintf(constants.FilterByBankAccountId, bankAccountId)
	}
	if billingType != "" {
		countQuery += fmt.Sprintf(constants.FilterByBillingType, billingType)
	}
	if studentId != 0 {
		countQuery += fmt.Sprintf(constants.FilterByStudentId, studentId)
	}

	// Hitung total data
	if err := configs.DB.Raw(countQuery).Scan(&totalData).Error; err != nil {
		fmt.Println("Error executing count query:", err)
		return nil, 0, 0, fmt.Errorf("failed to get total count: %v", err)
	}

	fmt.Println("totalData", countQuery)

	// Hitung total halaman
	if totalData > 0 {
		totalPages = int((totalData + int64(limit) - 1) / int64(limit))
	}

	// Query utama untuk mengambil data dengan pagination
	mainQuery := `
		SELECT  
			bs.detail_billing_name,
			b.billing_type AS billing_type,
			CONCAT(s.nis, ' - ', s.full_name) AS student_name,
			sg.school_grade_name AS school_grade_name,
			sc.school_class_name AS school_class_name,
			bs.amount AS amount,
			sy.school_year_name AS school_year_name,
			b.bank_account_id AS bank_account_id,
			ba.bank_name AS bank_name,
			ba.account_number AS account_number,
			CASE
				WHEN bs.payment_status = '1' THEN 'belum bayar'
				WHEN bs.payment_status = '2' THEN 'lunas'
				ELSE 'Status Tidak Diketahui'
			END AS payment_status
		FROM billing_students bs
		LEFT JOIN billings b ON b.id = bs.billing_id
		LEFT JOIN students s ON s.id = bs.student_id
		LEFT JOIN school_grades sg ON sg.id = s.school_grade_id
		LEFT JOIN school_years sy ON sy.id = b.school_year_id  
		LEFT JOIN school_classes sc ON sc.id = s.school_class_id
		LEFT JOIN bank_accounts ba ON ba.id = b.bank_account_id
		JOIN user_students us ON us.student_id = s.id
		JOIN user_schools usc ON usc.user_id = us.user_id
		WHERE bs.deleted_at IS NULL
	`
	if user.UserSchool != nil {
		mainQuery += fmt.Sprintf(constants.FilterByUscSchoolId, user.UserSchool.SchoolID)
	}
	if schoolClassId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterByUsSchoolClassId, schoolClassId)
	}
	if schoolGradeId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterBySgSchoolGradeId, schoolGradeId)
	}
	if paymentStatusId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterByBsPaymentStatus, paymentStatusId)
	}
	if schoolYearId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterBySySchoolYearId, schoolYearId)
	}
	if bankAccountId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterByBankAccountId, bankAccountId)
	}
	if billingType != "" {
		mainQuery += fmt.Sprintf(constants.FilterByBillingType, billingType)
	}
	if studentId != 0 {
		mainQuery += fmt.Sprintf(constants.FilterByStudentId, studentId)
	}

	if sortBy != "" {
		mainQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	} else {
		mainQuery += " ORDER BY CASE WHEN bs.updated_at IS NOT NULL THEN 0 ELSE 1 END, bs.updated_at DESC, bs.created_at DESC"
	}
	mainQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	fmt.Println("mainQuery", mainQuery)

	// Eksekusi query utama
	if err := configs.DB.Raw(mainQuery).Scan(&billingReports).Error; err != nil {
		fmt.Println("Error executing query:", err)
		return nil, 0, 0, fmt.Errorf("failed to get billing report: %v", err)
	}

	return billingReports, totalPages, totalData, nil
}

func (r *billingReportRepository) GetSummaryBillingReport(schoolGradeId int, schoolClassId int, paymentStatusId int, schoolYearId int,
	bankAccountId int, billingType string, studentId int, user models.User) (response.BillingReportSummary, error) {
	var summary response.BillingReportSummary

	query := `
		SELECT 
			SUM(bs.amount) as total_billing_amount,
			SUM(CASE WHEN bs.payment_status = '2' THEN bs.amount ELSE 0 END) as total_pay_amount,
			SUM(CASE WHEN bs.payment_status = '1' THEN bs.amount ELSE 0 END) as total_not_pay_amount,
			COUNT(CASE WHEN bs.payment_status = '2' THEN 1 END) as total_billing_pay,
			COUNT(CASE WHEN bs.payment_status = '1' THEN 1 END) as total_billing_not_pay,
			COUNT (distinct bs.student_id) as total_student
		FROM billing_students bs 
		LEFT JOIN billings b ON b.id = bs.billing_id
		LEFT JOIN students s ON s.id = bs.student_id
		LEFT JOIN school_grades sg ON sg.id = s.school_grade_id
		LEFT JOIN school_years sy ON sy.id = b.school_year_id  
		LEFT JOIN school_classes sc ON sc.id = s.school_class_id
		LEFT JOIN bank_accounts ba ON ba.id = b.bank_account_id
		JOIN user_students us ON us.student_id = s.id
		JOIN user_schools usc ON usc.user_id = us.user_id
		WHERE bs.deleted_at is null
	`

	if user.UserSchool != nil {
		query += fmt.Sprintf(constants.FilterByUscSchoolId, user.UserSchool.SchoolID)
	}
	if schoolClassId != 0 {
		query += fmt.Sprintf(constants.FilterByUsSchoolClassId, schoolClassId)
	}
	if schoolGradeId != 0 {
		query += fmt.Sprintf(constants.FilterBySgSchoolGradeId, schoolGradeId)
	}
	if paymentStatusId != 0 {
		query += fmt.Sprintf(constants.FilterByBsPaymentStatus, paymentStatusId)
	}
	if schoolYearId != 0 {
		query += fmt.Sprintf(constants.FilterBySySchoolYearId, schoolYearId)
	}
	if bankAccountId != 0 {
		query += fmt.Sprintf(constants.FilterByBankAccountId, bankAccountId)
	}
	if billingType != "" {
		query += fmt.Sprintf(constants.FilterByBillingType, billingType)
	}
	if studentId != 0 {
		query += fmt.Sprintf(constants.FilterByStudentId, studentId)
	}
	err := configs.DB.Raw(query).Scan(&summary).Error
	if err != nil {
		return response.BillingReportSummary{}, err
	}
	return summary, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
