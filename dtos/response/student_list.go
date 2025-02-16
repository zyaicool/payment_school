package response

import (
	"time"
)

type StudentListResponse struct {
	Page      int                     `json:"page"`
	Limit     int                     `json:"limit"`
	TotalPage int                     `json:"totalPage"`
	TotalData int64                   `json:"totalData"`
	Data      []DetailStudentResponse `json:"data"`
}

type DetailStudentResponse struct {
	ID                 uint       `json:"id"`
	Nisn               string     `json:"nisn"`
	RegistrationNumber string     `json:"registrationNumber"`
	Nis                string     `json:"nis"`
	Nik                string     `json:"nik"`
	FullName           string     `json:"fullName"`
	Gender             string     `json:"gender"`
	Religion           string     `json:"religion"`
	Citizenship        string     `json:"citizenship"`
	BirthPlace         string     `json:"birthPlace"`
	BirthDate          *time.Time `json:"birthDate"`
	Address            string     `json:"address"`
	SchoolGradeID      uint       `json:"schoolGradeId"`
	SchoolGrade        string     `json:"schoolGrade"`
	SchoolClassID      uint       `json:"schoolClassId"`
	SchoolClass        string     `json:"schoolClass"`
	SchoolID           uint       `json:"schoolId"`
	SchoolName         string     `json:"schoolName"`
	NoHandphone        string     `json:"noHandphone"`
	Height             string     `json:"height"`
	Weight             string     `json:"weight"`
	MedicalHistory     string     `json:"medicalHistory"`
	DistanceToSchool   uint       `json:"distanceToSchool"`
	Sibling            string     `json:"sibling"`
	NickName           string     `json:"nickName"`
	Email              string     `json:"email"`
	EntryYear          string     `json:"yearId"`
	Status             string     `json:"status"`
	Image              string     `json:"image"`
	Placeholder        string     `json:"placeholder"`
	SchoolYearID       uint       `json:"schoolYearId"`
	SchoolYearName     string     `json:"schoolYearName"`
}
