package request

import (
	"encoding/json"
	"fmt"
	"time"
)

type UpdateStudentRequest struct {
	Nisn             string `json:"nisn"`
	Nik              string `json:"nik"`
	Nis              string `json:"nis"`
	SchoolYearID     uint   `json:"schoolYearId"`
	SchoolGradeID    uint   `json:"schoolGradeId"`
	SchoolClassID    uint   `json:"schoolClassId"`
	FullName         string `json:"fullName" validate:"required"`
	Gender           string `json:"gender" validate:"required"`
	Religion         string `json:"religion" validate:"required"`
	Citizenship      string `json:"citizenship" validate:"required"`
	BirthPlace       string `json:"birthPlace" validate:"required"`
	BirthDate        string `json:"birthDate" validate:"required"`
	Address          string `json:"address" validate:"required"`
	NoHandphone      string `json:"noHandphone" validate:"required"`
	Height           string `json:"height"`
	Weight           string `json:"weight"`
	MedicalHistory   string `json:"medicalHistory"`
	DistanceToSchool uint   `json:"distanceToSchool" validate:"required"`
	Sibling          string `json:"sibling"`
	NickName         string `json:"nickName"`
	Image            string `json:"image"`
	Status           string `json:"status"`
	EmailParent      string `json:"emailParent" validate:"email"`
}

// Custom UnmarshalJSON for BirthDate
func (r *UpdateStudentRequest) UnmarshalJSON(data []byte) error {
	type Alias UpdateStudentRequest
	aux := &struct {
		BirthDate string `json:"birthDate"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Validate the date string format
	if aux.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.BirthDate)
		if err != nil {
			return fmt.Errorf("invalid birthDate format: %v", err)
		}
		// Jika parsing berhasil, set BirthDate sebagai time.Time
		r.BirthDate = parsedDate.Format("2006-01-02")
	}

	return nil
}
