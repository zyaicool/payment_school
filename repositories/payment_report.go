package repositories

import (
	"fmt"
	"math/big"
	database "schoolPayment/configs"
	"schoolPayment/dtos/response"
	"time"
)

// PaymentReportRepository defines the interface for payment report operations
type PaymentReportRepository interface {
	// GetPaymentReport retrieves paginated payment report details with optional filters
	GetPaymentReport(page int, limit int, sortBy string, sortOrder string, paymentTypeId int, userId int, startDate time.Time, endDate time.Time, studentId int, schoolID int, isAllData bool, paymentStatus string, schoolGradeId int, schoolClassId int) ([]response.PaymentReportDetail, int, int64, error)

	// GetPaymentReportSummary retrieves payment summary statistics with optional filters
	GetPaymentReportSummary(paymentTypeId int, userId int, startDate time.Time, endDate time.Time, studentId int, schoolID int, paymentStatus string, schoolGradeId int, schoolClassId int) (*response.PaymentReportResponse, error)
}

type paymentReportRepository struct{}

// NewPaymentReportRepository creates a new instance of PaymentReportRepository
func NewPaymentReportRepository() PaymentReportRepository {
	return &paymentReportRepository{}
}

// GetPaymentReport retrieves paginated payment report details with the following features:
// - Joins multiple tables to get comprehensive payment information
// - Supports filtering by school, student, payment type, date range and user
// - Implements pagination and sorting
// - Returns payment details, total pages and total records
func (r *paymentReportRepository) GetPaymentReport(
	page int, limit int,
	sortBy, sortOrder string,
	paymentTypeId, userId int,
	startDate, endDate time.Time,
	studentId, schoolID int,
	isAllData bool, paymentStatus string, schoolGradeId int, schoolClassId int,
) ([]response.PaymentReportDetail, int, int64, error) {
	var reports []response.PaymentReportDetail
	var total int64
	offset := (page - 1) * limit

	query := database.DB.Table("transaction_billings as tb").
		Select(`
			distinct
			tb.id,
			tb.invoice_number,
			CONCAT(s.nis,' - ',s.full_name) as student_name,
			tbd.created_at as payment_date,
			CASE 
			WHEN tb.transaction_type = 'PT01' THEN 'kasir'
				ELSE (
					SELECT concat(mp.payment_method, '-', mp.bank_name) 
					FROM master_payment_method mp
					WHERE mp.id = tbd.master_payment_method_id
					LIMIT 1
				)
			END AS payment_method,
			CASE 
            	WHEN u.role_id = 2 THEN '-' 
            	ELSE COALESCE((SELECT u.username FROM users u WHERE u.id = tb.created_by), 'system') 
       		END AS username, 
			sg.school_grade_name,
			sc.school_class_name,
			tb.total_amount,
			CASE 
				WHEN tb.transaction_status = 'PS01' THEN 'menunggu'
				WHEN tb.transaction_status = 'PS02' THEN 'lunas'
				WHEN tb.transaction_status = 'PS03' THEN 'gagal'
				ELSE tb.transaction_status  -- Default case, if no match
			end as transaction_status,
			tb.created_at
		`).
		Joins("JOIN transaction_billing_details tbd ON tbd.transaction_billing_id = tb.id").
		Joins("JOIN billing_students bs ON bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])").
		Joins("JOIN students s ON s.id = tb.student_id").
		Joins("JOIN school_grades sg ON sg.id = s.school_grade_id").
		Joins("JOIN school_classes sc ON sc.id = s.school_class_id").
		Joins("JOIN user_students us ON us.student_id = s.id").
		Joins("JOIN user_schools usc ON usc.user_id = us.user_id").
		Joins("JOIN users u ON u.id = usc.user_id").
		Where("tb.deleted_at IS NULL and tb.transaction_status != 'PS03'")

		// Get total count
	countQuery := database.DB.Table("transaction_billings as tb").
		Joins("JOIN transaction_billing_details tbd ON tbd.transaction_billing_id = tb.id").
		Joins("JOIN billing_students bs ON bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])").
		Joins("JOIN students s ON s.id = tb.student_id").
		Joins("JOIN school_grades sg ON sg.id = s.school_grade_id").
		Joins("JOIN school_classes sc ON sc.id = s.school_class_id").
		Joins("JOIN user_students us ON us.student_id = s.id").
		Joins("JOIN user_schools usc ON usc.user_id = us.user_id").
		Joins("JOIN users u ON u.id = usc.user_id").
		Where("tb.deleted_at IS NULL AND tb.transaction_status != 'PS03'").
		Select("COUNT(DISTINCT tb.id)")

	// Apply filters
	if schoolID != 0 {
		query = query.Where("usc.school_id = ?", schoolID)
		countQuery = countQuery.Where("usc.school_id = ?", schoolID)
	}

	if studentId != 0 {
		query = query.Where("tb.student_id = ?", studentId)
	}

	if paymentTypeId != 0 {
		if paymentTypeId == 99 {
			query = query.Where("tb.transaction_type = 'PT01'")
		} else {
			query = query.Where("tbd.master_payment_method_id = ?", paymentTypeId)
		}
	}

	if schoolClassId != 0 {
		query = query.Where("sc.id = ? ", schoolClassId)
	}

	if schoolGradeId != 0 {
		query = query.Where("sg.id = ? ", schoolGradeId)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("tb.created_at BETWEEN ? AND ?", startDate, endDate)
	}
	if userId != 0 {
		query = query.Where("tb.created_by = ?", userId)
	}
	if paymentStatus != "" {
		query = query.Where("tb.transaction_status = ?", paymentStatus)
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	// Apply sorting
	if sortBy != "" {
		// Add table alias to the sortBy column to avoid ambiguity
		if sortBy == "created_at" {
			sortBy = "tb.created_at"
		}
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("tb.created_at DESC")
	}

	// Apply pagination
	if !isAllData {
		query = query.Offset(offset).Limit(limit)
	}

	// Execute query
	err = query.Find(&reports).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := (int(total) + limit - 1) / limit

	return reports, totalPages, total, nil
}

// GetPaymentReportSummary retrieves payment summary statistics including:
// - Total billing amount
// - Total paid amount
// - Total unpaid amount
// - Count of paid/unpaid billings
// - Total number of students
// The summary can be filtered by:
// - Payment type (specific type or PT01 for special case)
// - User ID
// - Date range
// - Student ID
// - School ID
func (r *paymentReportRepository) GetPaymentReportSummary(
	paymentTypeId int,
	userId int,
	startDate time.Time,
	endDate time.Time,
	studentId int,
	schoolID int, paymentStatus string, schoolGradeId int, schoolClassId int,
) (*response.PaymentReportResponse, error) {
	var summary response.PaymentReportSummary

	query := `
		select
			COUNT(CASE WHEN tb.transaction_status = 'PS01' OR tb.transaction_status = 'PS02' THEN 1 END) as total_transaction,
			COUNT(DISTINCT tb.student_id) as total_student,
			SUM(tb.total_amount) as total_transaction_amount
		from
			transaction_billings as tb
		join transaction_billing_details tbd on tbd.transaction_billing_id = tb.id
		join billing_students bs on	bs.id = any(string_to_array(tb.billing_student_ids,	',')::int[])
		join students s on s.id = tb.student_id
		join school_grades sg on sg.id = s.school_grade_id
		join school_classes sc on sc.id = s.school_class_id
		join user_students us on us.student_id = s.id
		join user_schools usc on usc.user_id = us.user_id
		where tb.deleted_at is null
	`

	// Apply filters
	if schoolID != 0 {
		query += fmt.Sprintf(" AND usc.school_id = %d", schoolID)
	}

	if studentId != 0 {
		query += fmt.Sprintf(" AND s.id = %d", studentId)
	}

	if paymentTypeId != 0 {
		if paymentTypeId == 99 {
			query += " AND tb.transaction_type = 'PT01'"
		} else {
			query += fmt.Sprintf(" AND tbd.master_payment_method_id = %d", paymentTypeId)
		}
	}

	if schoolClassId != 0 {
		query += fmt.Sprintf(" AND sc.id = %d ", schoolClassId)
	}

	if schoolGradeId != 0 {
		query += fmt.Sprintf(" AND sg.id = %d ", schoolGradeId)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query += fmt.Sprintf(" AND tb.created_at BETWEEN '%s' AND '%s'", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))
	}

	if userId != 0 {
		query += fmt.Sprintf(" AND tb.created_by = %d", userId)
	}
	if paymentStatus != "" {
		query += fmt.Sprintf(" AND tb.transaction_status = '%s' ", paymentStatus)
	}

	err := database.DB.Raw(query).Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	return &response.PaymentReportResponse{
		TotalTransactionAmount: FormatBigInt(summary.TotalTransactionAmount),
		TotalTransaction:       summary.TotalTransaction,
		TotalStudent:           summary.TotalStudent,
		ListPaymentReport: response.ListPaymentReport{
			Page:      0,
			Limit:     0,
			TotalPage: 0,
			TotalData: 0,
			Data:      []response.PaymentReportDetail{},
		},
	}, nil
}

func FormatBigInt(value string) *big.Int {
	num := new(big.Int)
	rsp, _ := num.SetString(value, 10)
	return rsp
}
