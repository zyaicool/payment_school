package request

type SchoolCreateUpdateRequest struct {
	Npsn             int    `json:"npsn"`
	SchoolGradeID    int    `json:"schoolGradeId"`
	SchoolName       string `json:"schoolName"`
	SchoolProvince   string `json:"schoolProvince"`
	SchoolCity       string `json:"schoolCity"`
	SchoolPhone      string `json:"schoolPhone"`
	SchoolAddress    string `json:"schoolAddress"`
	SchoolMail       string `json:"schoolMail"`
	SchoolFax        string `json:"schoolFax"`
	SchoolLogo       string `json:"schoolLogo"`
	SchoolLetterhead string `json:"schoolLetterhead"`
}
