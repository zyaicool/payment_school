package models

type SchoolGrade struct {
	Master
	SchoolGradeCode string        `json:"schoolGradeCode"`
	SchoolGradeName string        `json:"schoolGradeName"`
	SchoolClass     []SchoolClass `gorm:"foreignKey:SchoolGradeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"schoolClass"`
}
