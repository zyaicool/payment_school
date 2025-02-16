package response

type DetailSchoolResponse struct {
	ID               int    `json:"id"`
	Npsn             int    `json:"npsn"`
	SchoolCode       string `json:"schoolCode"`
	SchoolName       string `json:"schoolName"`
	SchoolProvince   string `json:"schoolProvince"`
	SchoolCity       string `json:"schoolCity"`
	SchoolPhone      string `json:"schoolPhone"`
	SchoolAddress    string `json:"schoolAddress"`
	SchoolMail       string `json:"schoolMail"`
	SchoolFax        string `json:"schoolFax"`
	SchoolLogo       string `json:"schoolLogo"`
	SchoolGradeID    uint   `json:"schoolGradeId"`
	SchoolGradeName  string `json:"schoolGradeName"`
	SchoolLetterhead string `json:"schoolLetterhead"`
}

type SchoolOnBoardingResponse struct {
	ID         uint   `json:"id"`
	SchoolName string `json:"schoolName"`
	Npsn       int    `json:"npsn"`
	SchoolLogo string `json:"schoolLogo"`
}