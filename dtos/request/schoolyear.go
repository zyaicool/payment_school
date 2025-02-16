package request

type SchoolYearCreateUpdateRequest struct {
	SchoolYearName string `json:"schoolYearName"`
	SchoolId       int    `json:"schoolId"`
	StartDate      string `json:"startDate"` // Use string or time.Time depending on your parsing preference
	EndDate        string `json:"endDate"`   // Same as abo
}
