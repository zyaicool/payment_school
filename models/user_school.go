package models

type UserSchool struct {
	Master
	SchoolID uint    `json:"schoolId"`
	UserID   uint    `json:"userId"`
	School   *School `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
	User     *User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
}
