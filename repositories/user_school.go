package repositories

import (
	database "schoolPayment/configs"
	"schoolPayment/models"
)

func GetUserSchoolByUserId(id uint) (models.UserSchool, error) {
	var userSchool models.UserSchool
	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", id).First(&userSchool)
	return userSchool, result.Error
}

func CreateUserSchool(userSchool *models.UserSchool) (*models.UserSchool, error) {
	result := database.DB.Create(&userSchool)
	return userSchool, result.Error
}

func UpdateUserSchool(userSchool *models.UserSchool) (*models.UserSchool, error) {
	result := database.DB.Save(&userSchool)
	return userSchool, result.Error
}
