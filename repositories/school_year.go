package repositories

import (
	"fmt"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type SchoolYearRepository interface {
	GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, schoolId uint) ([]models.SchoolYearList, int64, int, error)
	CreateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error)
	UpdateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error)
	GetLastSequenceNumberSchoolYears() (int, error)
	GetAllSchoolYearsBulk(yearNames []string, schoolID uint) ([]models.SchoolYear, error)
	GetSchoolYearByID(id uint) (models.SchoolYear, error)
}

// Struct definition
type schoolYearRepository struct {
	db *gorm.DB
}

// Constructor function
func NewSchoolYearRepository(db *gorm.DB) SchoolYearRepository {
	return &schoolYearRepository{db: db}
}

func (r *schoolYearRepository) GetAllSchoolYear(page int, limit int, search string, sortBy string, sortOrder string, schoolId uint) ([]models.SchoolYearList, int64, int, error) {
	var schoolYearList []models.SchoolYearList
	var total int64
	var totalPages int

	query := r.db.Model(&models.SchoolYear{}).Select("school_years.*, users.username as create_by_username").Where("school_years.deleted_at IS NULL AND school_years.school_id = ?", schoolId)
	if search != "" {
		query = query.Where("school_year_name LIKE ?", "%"+search+"%")
	}

	query = query.Joins("JOIN users ON school_years.created_by = users.id")

	// Count total records without pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, fmt.Errorf("error counting total records: %w", err)
	}

	if limit > 0 {
		totalPages = int(total / int64(limit))
		if total%int64(limit) != 0 {
			totalPages++
		}
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if sortBy != "" {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("CASE WHEN school_years.updated_at IS NOT NULL THEN 0 ELSE 1 END, school_years.updated_at DESC, school_years.created_at DESC")
	}

	result := query.Find(&schoolYearList)
	if result.Error != nil {
		return nil, total, totalPages, fmt.Errorf("error fetching records: %w", result.Error)
	}

	return schoolYearList, total, totalPages, nil
}

func (r *schoolYearRepository) GetSchoolYearByID(id uint) (models.SchoolYear, error) {
	var schoolYear models.SchoolYear
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&schoolYear)
	return schoolYear, result.Error
}

func (schoolYearRepository *schoolYearRepository) CreateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error) {
	result := database.DB.Create(&schoolYear)
	return schoolYear, result.Error
}

func (schoolYearRepository *schoolYearRepository) UpdateSchoolYear(schoolYear *models.SchoolYear) (*models.SchoolYear, error) {
	result := database.DB.Save(&schoolYear)
	return schoolYear, result.Error
}

func (schoolYearRepository *schoolYearRepository) GetLastSequenceNumberSchoolYears() (int, error) {
	var lastSequence int
	result := database.DB.
		Table("school_years").
		Select("COALESCE(MAX(id), 0)").
		Scan(&lastSequence)

	if result.Error != nil {
		return 0, result.Error
	}

	return lastSequence, nil
}

func GetLatestSchoolYear() (models.SchoolYear, error) {
	var schoolYear models.SchoolYear
	query := database.DB.Where("deleted_at IS NULL")

	result := query.Order("created_at ASC").First(&schoolYear)
	return schoolYear, result.Error
}

func (r *schoolYearRepository) GetAllSchoolYearsBulk(yearNames []string, schoolID uint) ([]models.SchoolYear, error) {
	var years []models.SchoolYear
	result := r.db.Where("school_year_name IN ? AND school_id = ?", yearNames, schoolID).Find(&years)
	return years, result.Error
}
