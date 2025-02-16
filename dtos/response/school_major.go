package response

type SchoolMajorResponse struct {
	Data []DetailSchoolMajorResponse `json:"data"`
}

type DetailSchoolMajorResponse struct {
	ID              uint   `json:"id"`
	SchoolMajorName string `json:"schoolMajorName"`
}
