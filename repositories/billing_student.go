package repositories

import (
	"fmt"
	"strings"
	"time"

	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"

	"gorm.io/gorm"
	// "strconv"
)

type BillingStudentRepository interface {
	GetAllBillingStudent(page int, limit int, search string, studentId int, roleID int, schoolGradeID int, paymentType string, schoolID int, sortBy string, sortOrder string, schoolClassID int) ([]response.DataListBillingPerStudent, int, int64, error)
	GetDetailBillingStudentByID(billingStudentId int) (models.BillingStudent, error)
	UpdateBillingStudent(billingStudent *models.BillingStudent) error
	DeleteBillingStudent(billingStudentId int, deletedBy int) error
	GetTotalAmountByBillingStudentIds(billingStudentIds []string) (int, error)
	GetListBillingId(billingStudentIds []string) ([]int, error)
	// GetFirstBilling(studentId int) ([]response.LatestBillingStudent, error)
	// CreateBillingStudent(billingStudent *models.BillingStudent) (*models.BillingStudent, error)
	// GetDetailBillingIDRepositories(billingId int, studentId int) (response.BillingStudentByStudentIDBillingID, error)
	// GetBillingStudent(billingId int, studentId int) (*models.BillingStudent, error)
	// UpdateBillingStudent(billingStudent *models.BillingStudent) error
	CheckBillingStudentExists(studentID, billingID, billingDetailID uint) (bool, error)
	GetBillingDetailsByIDs(ids []uint) ([]*models.BillingDetail, error)
	CheckBulkBillingStudentExists(studentID uint, detailIDs []uint) ([]models.BillingStudent, error)
	BulkCreateBillingStudents(billingStudents []models.BillingStudent) ([]models.BillingStudent, error)
	GetInstallmentDetails(studentId int, schoolId uint) ([]response.BillingStudentByStudentIDDetailResponse, []response.DonationBillingResponse, error)
}

type billingStudentRepository struct{
	db *gorm.DB
}

func NewBillingStudentRepository(db *gorm.DB) BillingStudentRepository {
	return &billingStudentRepository{db: db}
}

func (billingStudentRepository *billingStudentRepository) GetAllBillingStudent(page int, limit int, search string, studentId int, roleID int, schoolGradeID int, paymentType string, schoolID int, sortBy string, sortOrder string, schoolClassID int) ([]response.DataListBillingPerStudent, int, int64, error) {
	var billingStudent []response.DataListBillingPerStudent
	var total int64
	totalPages := 0
	offset := (page - 1) * limit

	search = strings.ToLower(search)

	// Base query for getting data
	query := `
		select 
			bs.id as billing_student_id,
			bs.detail_billing_name,
			b.billing_type,
			sg.school_grade_name,
			sc.school_class_name,
			CONCAT(s.nis, ' - ', s.full_name) as student_name,
			bs.amount
			from billing_students bs 
			join billings b on b.id = bs.billing_id
			join students s on s.id = bs.student_id 
			join school_grades sg on sg.id = s.school_grade_id
			join school_classes sc on sc.id = s.school_class_id
			join user_students ust on ust.student_id = s.id
			join user_schools usc on usc.user_id = ust.user_id
			join schools sch on	sch.id = usc.school_id
			where bs.deleted_at is null 
			and bs.payment_status='1'
			and b.is_donation = false
	`

	// query += fmt.Sprintf("Where s.deleted_at is null and b.deleted_at is null and bs.payment_status='1'")

	// Add search filters
	if roleID == 2 {
		query += fmt.Sprintf(" AND s.id = %d ", studentId)
		if search != "" {
			query += fmt.Sprintf(" AND (LOWER(s.full_name) LIKE '%s' OR LOWER(sc.school_class_name) LIKE '%s' OR s.registration_number LIKE '%s') ",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
		}
	} else if roleID == 1 || roleID == 3 || roleID == 4 || roleID == 5 {
		if search != "" {
			query += fmt.Sprintf(" AND (LOWER(s.full_name) LIKE '%s' OR s.nis LIKE '%s') ",
				"%"+search+"%", "%"+search+"%")
		}

		if schoolGradeID != 0 {
			query += fmt.Sprintf(" AND sg.id = %d", schoolGradeID)
		}

		if schoolID != 0 {
			query += fmt.Sprintf(" AND usc.school_id = %d", schoolID)
		}

		if studentId != 0 {
			query += fmt.Sprintf(" AND s.id = %d", studentId)
		}

		if schoolClassID != 0 {
			query += fmt.Sprintf(" AND sc.id = %d", schoolClassID)
		}
	}

	if limit != 0 {
		// Get the total count of records without pagination (Offset/Limit)
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*) FROM (%s) AS total_query
		`, query)

		if err := database.DB.Raw(countQuery).Scan(&total).Error; err != nil {
			return nil, 0, 0, err
		}

		// Calculate total pages
		totalPages = int((total + int64(limit) - 1) / int64(limit)) // Ensure to round up for total pages

	}

	// Validate and apply sorting
	// if sortBy != "" && sort != "desc" {
	// 	sort = "asc" // Default sorting order
	// }
	// query += fmt.Sprintf(" ORDER BY b.updated_at %s", sort)

	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)
	} else {
		query += " ORDER BY CASE WHEN bs.updated_at IS NOT NULL THEN 0 ELSE 1 END, bs.updated_at DESC, bs.created_at DESC"
	}

	// Apply pagination to the original query
	query += fmt.Sprintf(" LIMIT %d OFFSET %d;", limit, offset)

	// Execute the query to get the billing records
	result := database.DB.Raw(query).Scan(&billingStudent)

	// Return billingStudent, totalPages, total, and error
	return billingStudent, totalPages, total, result.Error
}

func GetFirstBilling(studentId int) ([]response.LatestBillingStudent, error) {
	var billingStudent []response.LatestBillingStudent
	query := `
		select s.registration_number as nis, s.full_name as student_name, b.billing_number, 
		b.billing_name, sc.school_class_name, b.created_at as created_date, bs.due_date, bs.student_id, b.billing_amount,
		'' as semester, case when bs.payment_status = '1' then 'belum bayar' else 'lunas' end as billing_status, b.billing_amount, b.id AS billing_id
		from students s
		join billing_students bs on bs.student_id = s.id
		join billings b on	b.id = bs.billing_id
		left join school_classes sc on sc.id = b.school_class_id and b.school_class_id != 0	 
	`
	query += fmt.Sprintf(" where s.deleted_at is null and b.deleted_at is null and bs.payment_status = '1' and bs.student_id = %d ORDER BY due_date DESC LIMIT 2", studentId)
	result := database.DB.Raw(query).Find(&billingStudent)
	return billingStudent, result.Error
}

func CreateBillingStudent(billingStudent *models.BillingStudent, studentNis string, isMaster bool) (*models.BillingStudent, error) {
	var existingBillingStudent models.BillingStudent

	if isMaster {
		result := database.DB.Create(&billingStudent)
		return billingStudent, result.Error
	} else {
		err := database.DB.Where("billing_id = ? AND student_id = ? AND detail_billing_name LIKE ?", billingStudent.BillingID, billingStudent.StudentID, "%"+studentNis+"%").
			First(&existingBillingStudent).Error

		if err == nil {
			result := database.DB.Model(&existingBillingStudent).Updates(billingStudent)
			return &existingBillingStudent, result.Error
		} else if err.Error() == "record not found" {
			result := database.DB.Create(&billingStudent)
			return billingStudent, result.Error
		} else {
			return nil, err
		}
	}

}

func (billingStudentRepository *billingStudentRepository) GetInstallmentDetails(studentId int, schoolId uint) ([]response.BillingStudentByStudentIDDetailResponse, []response.DonationBillingResponse, error) {
	var installmentDetails []response.BillingStudentByStudentIDDetailResponse
	var donations []response.DonationBillingResponse

	// First, get the student's status
	var student models.Student
	if err := billingStudentRepository.db.Where("id = ?", studentId).First(&student).Error; err != nil {
		return nil, nil, err
	}

	// Execute installment query
	installmentQuery := `
		SELECT DISTINCT
			bs.id AS billing_student_id,
			bs.detail_billing_name,
			bs.amount,
			bs.due_date,
			CASE 
				WHEN bs.payment_status = '1' THEN 'belum bayar' 
				ELSE 'lunas' 
			END AS payment_status,
			b.billing_type,
			tb.transaction_status,
			bs.updated_at,
			bs.created_at
		FROM billing_students bs
		JOIN billings b ON b.id = bs.billing_id AND b.is_donation = false
		LEFT JOIN (
			SELECT *
			FROM transaction_billings tb
			ORDER BY created_at desc limit 1
		) tb ON
			bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])
		WHERE 
			bs.deleted_at is null and 
			bs.student_id = ? 
			AND (
				bs.payment_status = '1' AND tb.transaction_status IS NULL
				OR
				bs.payment_status = '1' AND tb.transaction_status NOT IN ('PS01','PS02')
			)
		ORDER BY 
			bs.updated_at DESC NULLS LAST,
			bs.created_at DESC;
	`

	err := billingStudentRepository.db.Raw(installmentQuery, studentId).Scan(&installmentDetails).Error
	if err != nil {
		return installmentDetails, nil, err
	}

	// Only execute donation query if student status is "aktif"
	if student.Status == "aktif" {
		donationQuery := `
			SELECT DISTINCT 
				b.id as billing_id, 
				b.billing_name,
				b.updated_at,
				b.created_at
			FROM billings b
			LEFT JOIN billing_students bs ON bs.billing_id = b.id AND bs.student_id = ?
			LEFT JOIN (
				SELECT * 
				FROM transaction_billings tb
				ORDER BY created_at desc limit 1
			) tb ON
				bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])
			left join user_schools us on us.user_id = b.created_by and us.school_id = ?
			LEFT JOIN students s ON s.id = bs.student_id
			WHERE b.is_donation = true 
			AND b.deleted_at IS NULL
			AND (
				bs.id IS NULL 
				OR (
					bs.payment_status = '1' 
					AND (
						tb.transaction_status IS NULL 
						OR tb.transaction_status NOT IN ('PS01','PS02')
					)
				)
			)
			AND us.school_id = ?
			ORDER BY 
				b.updated_at DESC NULLS LAST,
				b.created_at DESC`

		err = billingStudentRepository.db.Raw(donationQuery, studentId, schoolId, schoolId).Scan(&donations).Error
		if err != nil {
			return installmentDetails, nil, err
		}
	}

	return installmentDetails, donations, nil
}

func GetBillingStudent(billingId int, studentId int) (*models.BillingStudent, error) {
	var billingStudent *models.BillingStudent
	result := database.DB.Where("student_id =? and billing_id = ?", studentId, billingId).First(&billingStudent)
	return billingStudent, result.Error
}

func (billingStudentRepository *billingStudentRepository) GetDetailBillingStudentByID(billingStudentId int) (models.BillingStudent, error) {
	var billingStudentDetail models.BillingStudent
	result := database.DB.Where("id =? and deleted_at is null", billingStudentId).First(&billingStudentDetail)
	return billingStudentDetail, result.Error
}

func (billingStudentRepository *billingStudentRepository) UpdateBillingStudent(billingStudent *models.BillingStudent) error {
	result := database.DB.Save(&billingStudent)
	return result.Error
}

func (billingStudentRepository *billingStudentRepository) DeleteBillingStudent(billingStudentId int, deletedBy int) error {
	now := time.Now()
	result := database.DB.Model(&models.BillingStudent{}).
		Where("id = ? AND deleted_at IS NULL", billingStudentId).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"deleted_by": deletedBy,
		})
	if result.Error != nil {
		return fmt.Errorf("gagal menghapus data billing student: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("data billing student tidak ditemukan atau sudah dihapus")
	}
	return nil
}

func (billingStudentRepository *billingStudentRepository) GetTotalAmountByBillingStudentIds(billingStudentIds []string) (int, error) {
	var totalAmount int
	// Sum the amount where id is in the provided billingStudentIds array and deleted_at is null
	result := database.DB.Model(&models.BillingStudent{}).Where("id IN (?) AND deleted_at IS NULL", billingStudentIds).Select("SUM(amount)").Scan(&totalAmount)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalAmount, nil
}

func (billingStudentRepository *billingStudentRepository) GetListBillingId(billingStudentIds []string) ([]int, error) {
	var listBillingId []int
	result := database.DB.Model(&models.BillingStudent{}).
		Distinct("billing_id").
		Where("id IN (?) AND deleted_at IS NULL", billingStudentIds).
		Pluck("billing_id", &listBillingId)
	if result.Error != nil {
		return []int{}, result.Error
	}
	return listBillingId, nil
}

func (billingStudentRepository *billingStudentRepository) CheckBillingStudentExists(studentID, billingID, billingDetailID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&models.BillingStudent{}).
		Where("student_id = ? AND billing_id = ? AND billing_detail_id = ? AND deleted_at IS NULL",
			studentID, billingID, billingDetailID).
		Count(&count)

	return count > 0, result.Error
}

func (billingStudentRepository *billingStudentRepository) GetBillingDetailsByIDs(ids []uint) ([]*models.BillingDetail, error) {
	var details []*models.BillingDetail
	result := database.DB.Where("id IN ?", ids).Find(&details)
	return details, result.Error
}

func (billingStudentRepository *billingStudentRepository) CheckBulkBillingStudentExists(studentID uint, detailIDs []uint) ([]models.BillingStudent, error) {
	var existingRecords []models.BillingStudent
	result := database.DB.Where("student_id = ? AND billing_detail_id IN ?", studentID, detailIDs).
		Find(&existingRecords)
	return existingRecords, result.Error
}

func (billingStudentRepository *billingStudentRepository) BulkCreateBillingStudents(billingStudents []models.BillingStudent) ([]models.BillingStudent, error) {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&billingStudents).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return billingStudents, nil
}

func GetBillingStudentsDashboard(studentID int, limit int) ([]response.BillingStudentByStudentIDDetailResponse, error) {
	var billingDetails []response.BillingStudentByStudentIDDetailResponse

	query := `
		SELECT DISTINCT
			bs.id AS billing_student_id,
			bs.detail_billing_name,
			bs.amount,
			bs.due_date,
			CASE 
				WHEN bs.payment_status = '1' THEN 'belum bayar' 
				ELSE 'lunas' 
			END AS payment_status,
			b.billing_type,
			tb.transaction_status
		FROM billing_students bs
		JOIN billings b ON b.id = bs.billing_id
		LEFT JOIN (
			SELECT * 
			FROM transaction_billings tb
			ORDER BY created_at desc limit 1
		) tb ON
			bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])
		WHERE 
			bs.student_id = ?
			AND (
				bs.payment_status = '1' AND tb.transaction_status IS NULL
				OR
				bs.payment_status = '1' AND tb.transaction_status NOT IN ('PS01', 'PS02', 'PS03')
			)
		LIMIT ?;
	`

	// Execute the query
	err := database.DB.Raw(query, studentID, limit).Scan(&billingDetails).Error
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return billingDetails, nil
}
