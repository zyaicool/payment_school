package repositories

import (
	"strings"

	database "schoolPayment/configs"
	"schoolPayment/constants"
	models "schoolPayment/models"
)

type StudentParentRepository struct{}

func NewStudentParentRepository() StudentParentRepository {
	return StudentParentRepository{}
}

func (studentParentRepository *StudentParentRepository) GetAllParent(page int, limit int, search string, user models.User) ([]models.StudentParent, error) {
	var parents []models.StudentParent
	query := database.DB.Where("student_parents.deleted_at IS NULL")
	if search != "" {
		query = query.Where("LOWER(student_parents.parent_name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if user.RoleID == 2 {
		query = query.Joins(constants.JoinUsersToStudentParenstAndFilterDeletedAt).
			Where("users.id = ?", user.ID)
	} else if user.RoleID == 5 {
		query = query.Joins(constants.JoinUsersToStudentParenstAndFilterDeletedAt).
			Joins(constants.JoinUserSchoolsToUsersAndFilterDeletedAt).
			Joins(constants.JoinSchoolsToUserSChoolsAndFilterDeletedAt).
			Where("schools.id = ?", user.UserSchool.School.ID)
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&parents)
	return parents, result.Error
}

func (studentParentRepository *StudentParentRepository) GetParentByID(id uint, user models.User) (models.StudentParent, error) {
	var parent models.StudentParent
	query := database.DB.Where("student_parents.id = ? AND student_parents.deleted_at IS NULL", id)
	if user.RoleID == 2 {
		query = query.Joins(constants.JoinUsersToStudentParenstAndFilterDeletedAt).
			Where("users.id = ?", user.ID)
	} else if user.RoleID == 5 {
		query = query.Joins(constants.JoinUsersToStudentParenstAndFilterDeletedAt).
			Joins(constants.JoinUserSchoolsToUsersAndFilterDeletedAt).
			Joins(constants.JoinSchoolsToUserSChoolsAndFilterDeletedAt).
			Where("schools.id = ?", user.UserSchool.School.ID)
	}
	result := query.First(&parent)
	return parent, result.Error
}

func (studentParentRepository *StudentParentRepository) CreateParent(parent *models.StudentParent) (*models.StudentParent, error) {
	result := database.DB.Create(&parent)
	return parent, result.Error
}

func (studentParentRepository *StudentParentRepository) UpdateParent(parent *models.StudentParent) (*models.StudentParent, error) {
	result := database.DB.Save(parent)

	return parent, result.Error
}

func GetParentByUserLogin(id uint) (*models.StudentParent, error) {
	var parent *models.StudentParent
	result := database.DB.Where("user_id = ? AND deleted_at IS NULL", id).First(&parent)
	return parent, result.Error
}
