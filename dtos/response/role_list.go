package response

import (
	models "schoolPayment/models"
)

type RoleListResponse struct {
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Data  []models.Role `json:"data"`
}

type DetailRole struct {
	ID         int                `json:"id"`
	RoleName   string             `json:"roleName"`
	RoleMatrix []DetailRoleMatrix `json:"roleMatrix"`
}

type DetailRoleMatrix struct {
	PageName string `json:"pageName"`
	PageCode string `json:"pageCode"`
	IsCreate bool   `json:"isCreate"`
	IsRead   bool   `json:"isRead"`
	IsUpdate bool   `json:"isUpdate"`
	IsDelete bool   `json:"isDelete"`
}
