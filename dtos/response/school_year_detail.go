package response

import "time"

type SchoolYearDetailResponse struct {
	ID             int        `json:"id"`
	SchoolYearName string     `json:"schoolYearName"`
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}
