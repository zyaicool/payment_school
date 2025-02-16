package repositories

import (
	"strings"

	database "schoolPayment/configs"
	models "schoolPayment/models"
)

type RoleRepository interface {
	GetAllRole(page int, limit int, search string, roleID int) ([]models.Role, error)
	GetRolesByNames(roleNames []string) (map[string]uint, error)
}

type roleRepository struct{}

func NewRoleRepository() RoleRepository {
	return &roleRepository{}
}

func (roleRepository *roleRepository) GetAllRole(page int, limit int, search string, roleID int) ([]models.Role, error) {
	var roles []models.Role
	query := database.DB.Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("lower(name) like ?", "%"+strings.ToLower(search)+"%")
	}

	if roleID == 5 {
		query = query.Where("id in (2,3,4,5)")
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&roles)
	return roles, result.Error
}

func GetRoleByID(id uint) (models.Role, error) {
	var role models.Role
	result := database.DB.Where("id = ? AND deleted_at IS NULL", id).Preload("RoleMatrix").First(&role)
	return role, result.Error
}

func GetRoleNames() ([]string, error) {
	var names []string
	result := database.DB.Model(&models.Role{}).
		Where("deleted_at IS NULL").
		Select("name").
		Order("name ASC"). // Optional: for sorted results
		Pluck("name", &names)

	return names, result.Error
}

func (roleRepository *roleRepository) GetRolesByNames(roleNames []string) (map[string]uint, error) {
	var roles []models.Role
	roleMap := make(map[string]uint)

	err := database.DB.Where("name IN ? AND deleted_at IS NULL", roleNames).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	return roleMap, nil
}
