package models

type UserStudent struct {
	Master
	UserID    uint    `json:"userId"`
	StudentID uint    `json:"studentId"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Student   Student `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"student"`
}
