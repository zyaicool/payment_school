package models

type User struct {
	Master
	RoleID         uint           `json:"roleId"`
	Username       string         `json:"username"`
	Role           Role           `gorm:"foreignKey:RoleID"`
	Email          string         `json:"email" validate:"required"`
	Password       string         `json:"password" validate:"required, min:8"`
	IsVerification bool           `json:"isVerification" gorm:"default:false"` // Set default to false
	IsBlock        bool           `json:"isBlock" gorm:"default:false"`        // Set default to false
	UserSchool     *UserSchool    `gorm:"foreignKey:UserID" json:"userSchool"`
	UserStudents   []*UserStudent `gorm:"foreignKey:UserID" json:"userStudents"`
	Image          string         `json:"image"`
}
