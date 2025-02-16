package models

type SchoolClass struct {
	Master
	SchoolID        uint         `json:"schoolId"`
	SchoolGradeID   uint         `json:"schoolGradeId"`
	SchoolClassCode string       `json:"schoolClassCode"`
	SchoolClassName string       `json:"schoolClassName"`
	PrefixClassID   int          `json:"prefixClassId"`
	SchoolMajorID   int          `json:"schoolMajorId"`
	Suffix          string       `json:"suffix"`
	SchoolGrade     *SchoolGrade `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"schoolGrade"`
	School          *School      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
}
