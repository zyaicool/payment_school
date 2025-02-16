package response

import (
	"time"
)

type SchoolYearListResponse struct {
	Data      []DetailSchoolYearResponse `json:"data"`
	TotalData int                        `json:"totalData"`
	TotalPage int                        `json:"totalPage"`
	Limit     int                        `json:"limit"`
	Page      int                        `json:"page"`
}

type DetailSchoolYearResponse struct {
	ID             uint       `json:"id"`
	SchoolYearName string     `json:"schoolYearName"`
	StartDate      *time.Time `json:"startDate"`
	EndDate        *time.Time `json:"endDate"`
	CreatedAt      time.Time  `json:"createdAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
