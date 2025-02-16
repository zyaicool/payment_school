package models

type StudentGuardian struct {
	Master
	UserID              uint   `json:"userId"`
	GuardianName        string `json:"guardianName"`
	GuardianAddress     string `json:"guardianAddress"`
	GuardianHandphone   int    `json:"guardianHandphone"`
	GuardianMail        string `json:"guardianMail"`
	GuardianCitizenship string `json:"guardianCitizenship"`
	GuardianSalary      string `json:"guardianSalary"`
}
