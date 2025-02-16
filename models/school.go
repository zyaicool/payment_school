package models

type School struct {
	Master
	ID               uint          `gorm:"primaryKey" json:"id"` // Primary key field (ID)
	Npsn             int           `json:"npsn"`
	SchoolCode       string        `json:"schoolCode"`
	SchoolName       string        `json:"schoolName"`
	SchoolProvince   string        `json:"schoolProvince"`
	SchoolCity       string        `json:"schoolCity"`
	SchoolPhone      string        `json:"schoolPhone"`
	SchoolAddress    string        `json:"schoolAddress"`
	SchoolMail       string        `json:"schoolMail"`
	SchoolFax        string        `json:"schoolFax"`
	SchoolLogo       string        `json:"schoolLogo"`
	SchoolLetterhead string        `json:"schoolLetterhead"`
	SchoolGradeID    uint          `json:"schoolGradeId"`
	SchoolGrade      *SchoolGrade  `gorm:"foreignKey:SchoolGradeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"schoolGrades"`
	UserSchools      []*UserSchool `gorm:"foreignKey:SchoolID" json:"userSchools"`
}

type SchoolList struct {
	School                   // Embeds the original Billing struct
	CreatedByUsername string `gorm:"column:created_by_username"`
}
