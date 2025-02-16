package models

import "time"

type Student struct {
	Master
	Nisn               string       `json:"nisn"`
	RegistrationNumber string       `json:"registrationNumber"`
	Nis                string       `json:"nis"`
	Nik                string       `json:"nik"`
	FullName           string       `json:"fullName" validate:"required"`
	Gender             string       `json:"gender" validate:"required"`
	Religion           string       `json:"religion" validate:"required"`
	Citizenship        string       `json:"citizenship" validate:"required"`
	BirthPlace         string       `json:"birthPlace" validate:"required"`
	BirthDate          *time.Time   `json:"birthDate" validate:"required"`
	Address            string       `json:"address" validate:"required"`
	SchoolGrade        string       `json:"schoolGrade" validate:"required"`
	SchoolGradeID      uint         `json:"schoolGradeId"`
	SchoolClass        string       `json:"schoolClass"`
	SchoolClassID      uint         `json:"schoolClassId"`
	NoHandphone        string       `json:"noHandphone" validate:"required"`
	Height             string       `json:"height"`
	Weight             string       `json:"weight"`
	MedicalHistory     string       `json:"medicalHistory"`
	DistanceToSchool   uint         `json:"distanceToSchool" validate:"required"`
	Sibling            string       `json:"sibling"`
	NickName           string       `json:"nickName"`
	Email              string       `json:"email" validate:"email"`
	EntryYear          string       `json:"yearId" validate:"required"`
	Status             string       `json:"status" validate:"required"`
	Image              string       `json:"image"`
	UserStudents       *UserStudent `gorm:"foreignKey:StudentID" json:"userStudents"`
	SchoolYearID       uint         `json:"schoolYearId"`
}
