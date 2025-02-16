package repositories

import (
	database "schoolPayment/configs"
	models "schoolPayment/models"
)

func GetAllGuardian(page int, limit int, search string) ([]models.StudentGuardian, error) {
	var guardians []models.StudentGuardian
	query := database.DB.Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("guardian_name like ?", "%"+search+"%")
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&guardians)
	return guardians, result.Error
}

func GetGuardianByID(id uint) (models.StudentGuardian, error) {
	var guardian models.StudentGuardian
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&guardian)
	return guardian, result.Error
}

func CreateGuardian(guardian *models.StudentGuardian) (*models.StudentGuardian, error) {
	result := database.DB.Create(&guardian)
	return guardian, result.Error
}

func UpdateGuardian(guardian *models.StudentGuardian) (*models.StudentGuardian, error) {
	result := database.DB.Save(guardian)

	return guardian, result.Error
}

func GetGuardianByUserLogin(id uint) (*models.StudentGuardian, error) {
	var guardian *models.StudentGuardian
	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", id).First(&guardian)
	return guardian, result.Error
}
