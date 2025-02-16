package response

type SchoolResponse struct {
	ID             int    `json:"id"`
	Npsn           int    `json:"npsn"`
	SchoolName     string `json:"schoolName"`
	SchoolProvince string `json:"schoolProvince"`
	SchoolPhone    string `json:"schoolPhone"`
	CreatedBy      string `json:"createdBy"`
}

type SchoolListResponse struct {
	Page      int              `json:"page"`
	Limit     int              `json:"limit"`
	TotalPage int              `json:"totalPage"`
	TotalData int              `json:"totalData"`
	Data      []SchoolResponse `json:"data"`
}
