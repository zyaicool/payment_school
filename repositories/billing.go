package repositories

import (
	"fmt"
	"schoolPayment/configs"
	database "schoolPayment/configs"
	response "schoolPayment/dtos/response"
	"schoolPayment/models"

	"gorm.io/gorm"
)


type BillingRepositoryInterface interface {
	GetDetailBillingsByBillingID(billingID uint) ([]response.DetailBilling, error)
	GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, user models.User) ([]models.BillingList, int, int64, error)
	GetBillingByID(id int) (models.Billing, error)
	CheckBillingCode(billingCode string) bool 
	CreateBillingDonation(billing *models.Billing) (*models.Billing, error)
	GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) ([]models.BillingStudentsExists, error)
	CheckBillingStudentExists(studentID uint, billingDetailID uint) (bool, error)
	CreateBilling(billing *models.Billing) (*models.Billing, error)
	GetLastSequenceNumberBilling() (int, error)
}

type BillingRepository struct{
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepositoryInterface {
	return &BillingRepository{db: db}
}

func (billingRepository BillingRepository) GetAllBilling(page int, limit int, search, billingType, paymentType, schoolGrade, sort string, sortBy string, sortOrder string, bankAccountId int, isDonation *bool, user models.User) ([]models.BillingList, int, int64, error) {
	var billingList []models.BillingList
	var total int64
	totalPages := 0

	// Base query to get the billing records
	query := billingRepository.db.Model(&models.Billing{}).
		Select("billings.*, CASE WHEN billings.is_donation = true THEN 'admin' ELSE users.username END as create_by_username").
		Joins("LEFT JOIN users ON users.id = billings.created_by").
		Joins("LEFT JOIN bank_accounts ON bank_accounts.id  = billings.bank_account_id ").
		Joins("LEFT JOIN school_grades ON school_grades.id = billings.school_grade_id").
		Joins("LEFT JOIN user_schools ON users.id = user_schools.user_id").
		Where("billings.deleted_at IS NULL")

	// Apply isDonation filter before executing the query
	if isDonation != nil {
		query = query.Where("is_donation = ?", *isDonation)
	}

	// Apply search filter
	if search != "" {
		query = query.Where("billing_name ILIKE ?", "%"+search+"%")
	}

	// Apply billing type filter
	if billingType != "" {
		query = query.Where("billing_type = ?", billingType)
	}

	// Apply payment type filter
	if paymentType != "" {
		query = query.Where("payment_type = ?", paymentType)
	}

	// Apply school grade filter
	if schoolGrade != "" {
		query = query.Where("school_grade_id = ?", schoolGrade)
	}

	// Apply account number filter
	if bankAccountId != 0 {
		query = query.Where("billings.bank_account_id = ?", bankAccountId)
	}

	fmt.Println("user school", user.UserSchool.SchoolID)
	if user.UserSchool.SchoolID != 0 {
		query = query.Where("user_schools.school_id = ?", user.UserSchool.SchoolID)
	}

	// Determine sort order
	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("created_at " + sort)
	}

	// Count total records without pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if limit != 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	// Execute the final query
	result := query.Preload("BankAccount").Preload("SchoolGrade").Find(&billingList)
	if result.Error != nil {
		return nil, 0, total, result.Error
	}

	return billingList, totalPages, total, nil
}

func (billingRepository *BillingRepository) GetBillingByID(id int) (models.Billing, error) {
	var billing models.Billing
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("SchoolYear").Preload("SchoolGrade").Preload("BankAccount").First(&billing)
	return billing, result.Error
}

func (billingRepository *BillingRepository) GetDetailBillingsByBillingID(billingID uint) ([]response.DetailBilling, error) {
	var detailBillings []response.DetailBilling
	query := database.DB.Model(&models.BillingDetail{}).
		Select("id, detail_billing_name, due_date, amount").Where("billing_id = ?", billingID)
	result := query.Find(&detailBillings)
	return detailBillings, result.Error
}

func (billingRepository BillingRepository) CreateBilling(billing *models.Billing) (*models.Billing, error) {
	result := billingRepository.db.Create(&billing)
	return billing, result.Error
}

func UpdateBilling(billing *models.Billing) (*models.Billing, error) {
	result := database.DB.Save(&billing)
	return billing, result.Error
}

func (billingRepository BillingRepository) GetLastSequenceNumberBilling() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("billings").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

func (billingRepository *BillingRepository) CheckBillingCode(billingCode string) bool {
	var billing models.Billing
	result := database.DB.Where("billing_code = ? AND deleted_at IS NULL", billingCode).First(&billing)
	if result.Error == nil {
		return false
	}

	return true
}

func (billingRepository *BillingRepository) CreateBillingDonation(billing *models.Billing) (*models.Billing, error) {
	billing.IsDonation = true
	result := configs.DB.Create(&billing)
	return billing, result.Error
}

func (billingRepository *BillingRepository) GetBillingByStudentID(studentID, schoolYearID, schoolGradeID, schoolClassID int) ([]models.BillingStudentsExists, error) {
	var billings []models.BillingStudentsExists

	// Main query building
	query := database.DB.Model(&models.Billing{}).
		Select(`
			billings.*,
			billing_details.id as billing_detail_id,
			billing_details.detail_billing_name,
			billing_details.amount,
			EXISTS(
				SELECT 1 FROM billing_students bs 
				WHERE bs.student_id = ? AND 
					  bs.billing_detail_id = billing_details.id AND 
					  bs.deleted_at IS NULL
			) as is_exist,
			 CASE 
				WHEN billing_details.due_date < CURRENT_DATE THEN true
				ELSE false
			END AS disabled
		`, studentID).
		Joins("JOIN billing_details ON billing_details.billing_id = billings.id").
		Where("billings.deleted_at IS NULL")

	// Apply filters based on parameters
	if schoolYearID > 0 {
		query = query.Where("billings.school_year_id = ?", schoolYearID)
	}

	if schoolGradeID > 0 {
		query = query.Where("billings.school_grade_id = ?", schoolGradeID)
	}

	if schoolClassID > 0 {
		query = query.Where("CONCAT(',', billings.school_class_ids, ',') LIKE ?",
			fmt.Sprintf("%%,%d,%%", schoolClassID))
	}

	// Execute query
	err := query.Find(&billings).Error
	if err != nil {
		return nil, err
	}

	return billings, nil
}

func (billingRepository *BillingRepository) CheckBillingStudentExists(studentID uint, billingDetailID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.BillingStudent{}).
		Where("student_id = ? AND billing_detail_id = ? AND deleted_at IS NULL",
			studentID, billingDetailID).
		Count(&count).Error

	return count > 0, err
}

func GetLatestDonation(isDonation bool, studentId uint, schoolId uint) (models.Billing, error) {
	var billing models.Billing

	donationQuery := `
        SELECT DISTINCT 
            b.*
        FROM billings b
        LEFT JOIN billing_students bs ON bs.billing_id = b.id AND bs.student_id = ?
        LEFT JOIN (
			SELECT * 
			FROM transaction_billings tb
			ORDER BY created_at desc limit 1
		) tb ON
			bs.id = ANY(string_to_array(tb.billing_student_ids, ',')::int[])
		left join user_schools us on us.user_id = b.created_by and us.school_id = ?
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
            b.created_at DESC LIMIT 1`

	err := database.DB.Raw(donationQuery, studentId, schoolId, schoolId).Scan(&billing).Error
	if err != nil {
		return models.Billing{}, err
	}

	return billing, err
}
