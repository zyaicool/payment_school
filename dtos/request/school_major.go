package request

type SchoolMajorCreate struct {
	SchoolMajorName string `json:"schoolMajorName"`
	SchoolID        uint   `json:"schoolId"`
}
