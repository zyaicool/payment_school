package models

type RoleMatrix struct {
	Master
	RoleID   uint   `json:"roleId"`
	PageName string `json:"pageName"`
	PageCode string `json:"pageCode"`
	IsCreate bool   `json:"isCreate"`
	IsRead   bool   `json:"isRead"`
	IsUpdate bool   `json:"isUpdate"`
	IsDelete bool   `json:"isDelete"`
}
