package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/models"

	"gorm.io/gorm"
)

type SchoolMajorRepository interface{
	GetAllSchoolMajorRepository(search string, user models.User) ([]models.SchoolMajor, error)
    CreateSchoolMajorRepository(major *models.SchoolMajor) (*models.SchoolMajor, error)
    CheckSchoolMajorExists(majorName string, schoolID uint) (bool, error)
}

type schoolMajorRepository struct{
	db *gorm.DB
}

func NewSchoolMajorRepository(db *gorm.DB) SchoolMajorRepository {
	return &schoolMajorRepository{db}
}

func (r *schoolMajorRepository) GetAllSchoolMajorRepository(search string, user models.User) ([]models.SchoolMajor, error) {
	var major []models.SchoolMajor

	query := database.DB.Where("school_majors.deleted_at IS NULL")

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(school_majors.school_major_name) LIKE ?", searchTerm)
	}

	if user.UserSchool != nil {
		query = query.Where("school_majors.school_id = ?", user.UserSchool.SchoolID)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&major)
	if result.Error != nil {
		return nil, result.Error
	}

	return major, nil
}

func (r *schoolMajorRepository)CreateSchoolMajorRepository(major *models.SchoolMajor) (*models.SchoolMajor, error) {
	result := database.DB.Create(&major)
	return major, result.Error
}

func (r *schoolMajorRepository)CheckSchoolMajorExists(majorName string, schoolID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.SchoolMajor{}).
		Where("LOWER(school_major_name) = ? AND school_id = ?", strings.ToLower(majorName), schoolID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
