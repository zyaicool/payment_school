// File: models/role.go
package models

type Role struct {
	Master
	Name       string       `json:"name" validate:"required"`
	RoleMatrix []RoleMatrix `gorm:"foreignKey:RoleID"`
}
