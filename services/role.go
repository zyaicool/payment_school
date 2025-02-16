package services

import (
	response "schoolPayment/dtos/response"
	models "schoolPayment/models"
	repositories "schoolPayment/repositories"
)

type RoleService struct {
	roleRepository repositories.RoleRepository
}

func NewRoleService(roleRepository repositories.RoleRepository) RoleService {
	return RoleService{roleRepository: roleRepository}
}

func (roleService *RoleService) GetAllRoles(page int, limit int, search string, roleID int) (response.RoleListResponse, error) {
	var mapRoles response.RoleListResponse
	mapRoles.Limit = limit
	mapRoles.Page = page

	// Call the repository to fetch the roles
	listRole, err := roleService.roleRepository.GetAllRole(page, limit, search, roleID)
	if err != nil {
		// If error occurs, return empty response with Data set to nil
		mapRoles.Data = nil
		return mapRoles, err
	}

	// If roles exist, return them, else return empty slice
	if len(listRole) > 0 {
		mapRoles.Data = listRole
	} else {
		mapRoles.Data = []models.Role{}
	}

	return mapRoles, nil
}


func GetRoleByID(id uint) (models.Role, error) {
	return repositories.GetRoleByID(id)
}
