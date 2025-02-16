package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/models"
)

type IPrefixClassRepository interface {
	GetAllPrefixClassRepository(search string, user models.User) ([]models.PrefixClass, error)
	CreatePrefixClassRepository(prefix *models.PrefixClass) (*models.PrefixClass, error)
	CheckPrefixClassExists(prefixName string, schoolID uint) (bool, error)
}

type PrefixClassRepository struct{}

func NewPrefixClassRepository() IPrefixClassRepository {
	return &PrefixClassRepository{}
}

func (r *PrefixClassRepository) GetAllPrefixClassRepository(search string, user models.User) ([]models.PrefixClass, error) {
	var prefix []models.PrefixClass

	query := database.DB.Where("prefix_classes.deleted_at IS NULL")

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(prefix_name) like ?", searchTerm)
	}

	if user.UserSchool != nil {
		query = query.Where("prefix_classes.school_id = ?", user.UserSchool.SchoolID)
	}

	query = query.Order("created_at DESC")

	result := query.Find(&prefix)
	if result.Error != nil {
		return nil, result.Error
	}

	return prefix, nil
}

func (r *PrefixClassRepository) CreatePrefixClassRepository(prefix *models.PrefixClass) (*models.PrefixClass, error) {
	result := database.DB.Create(&prefix)
	return prefix, result.Error
}

func (r *PrefixClassRepository) CheckPrefixClassExists(prefixName string, schoolID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.PrefixClass{}).
		Where("LOWER(prefix_name) = ? AND school_id = ?", strings.ToLower(prefixName), schoolID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
