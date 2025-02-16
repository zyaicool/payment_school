package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/models"
)

type BillingTypeRepository interface {
	GetAllBillingType(page int, limit int, search string) ([]models.BillingType, error)
	GetBillingTypeByIDTest(id uint) (models.BillingType, error)
	UpdateBillingType(billingType *models.BillingType) (*models.BillingType, error)
	GetLastSequenceNumberBillingType() (int, error)
	GetUserSchoolByUserIdTest(id uint) (models.UserSchool, error)
	CreateBillingType(billingType *models.BillingType) (*models.BillingType, error)
}

type billingTypeRepository struct{}

func NewBillingTypeRepository() BillingTypeRepository {
	return &billingTypeRepository{}
}

func (r *billingTypeRepository) GetAllBillingType(page int, limit int, search string) ([]models.BillingType, error) {
	var billingTypeList []models.BillingType
	query := database.DB.Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("LOWER(billing_type_name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if limit != 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	result := query.Find(&billingTypeList)
	return billingTypeList, result.Error
}

func GetBillingTypeByID(id uint) (models.BillingType, error) {
	var billingType models.BillingType
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("School").First(&billingType)
	return billingType, result.Error
}

func (r *billingTypeRepository) CreateBillingType(billingType *models.BillingType) (*models.BillingType, error) {
	result := database.DB.Create(&billingType)
	return billingType, result.Error
}

func (r *billingTypeRepository) UpdateBillingType(billingType *models.BillingType) (*models.BillingType, error) {
	result := database.DB.Save(&billingType)
	return billingType, result.Error
}

func (r *billingTypeRepository) GetLastSequenceNumberBillingType() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("billing_types").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

// Inii Hanyaa Sementaraa

func (r *billingTypeRepository) GetBillingTypeByIDTest(id uint) (models.BillingType, error) {
	var billingType models.BillingType
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("School").First(&billingType)
	return billingType, result.Error
}

func (r *billingTypeRepository) GetUserSchoolByUserIdTest(id uint) (models.UserSchool, error) {
	var userSchool models.UserSchool
	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", id).First(&userSchool)
	return userSchool, result.Error
}
