package response

type SchoolClassResponse struct {
	ID              int    `json:"id"`
	Unit            string `json:"unit"`
	PrefixClass     string `json:"prefixClass"`
	SchoolMajor     string `json:"schoolMajor"`
	SchoolClassName string `json:"schoolClassName"`
	CreatedBy       string `json:"createdBy"`
	Status          bool   `json:"status"`
	IsEdit          bool   `json:"isEdit"`
	Placeholder     string `json:"placeholder"`
}

type SchoolClassResponseRepo struct {
	ID                int    `json:"id"`
	Unit              string `json:"unit"`
	PrefixClass       string `json:"prefixClass"`
	SchoolMajor       string `json:"schoolMajor"`
	SchoolClassName   string `json:"schoolClassName"`
	CreatedBy         string `json:"createdBy"`
	Status            bool   `json:"status"`
	IsEdit            bool   `json:"isEdit"`
	Placeholder       string `json:"placeholder"`
	CreatedByUsername string `json:"createdByUsername"`
}

type SchoolClassListResponse struct {
	Page      int                   `json:"page"`
	Limit     int                   `json:"limit"`
	TotalPage int                   `json:"totalPage"`
	TotalData int64                 `json:"totalData"`
	Data      []SchoolClassResponse `json:"data"`
}
